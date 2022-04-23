package xb

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/mastercactapus/embedded/driver"
	"github.com/mastercactapus/embedded/serial/i2c"
	"github.com/mastercactapus/embedded/serial/spi"
)

type Client struct {
	r *bufio.Reader
	w io.Writer

	pinCount int
}

var (
	_ spi.ReadController      = (*spiClient)(nil)
	_ spi.ReadWriteController = (*spiClient)(nil)
	_ spi.WriteController     = (*spiClient)(nil)
)

type bw struct {
	io.Writer
}

func (bw bw) Write(p []byte) (int, error) {
	for i, b := range p {
		_, err := bw.Writer.Write([]byte{b})
		if err != nil {
			return i, err
		}
		time.Sleep(time.Millisecond)
	}
	return len(p), nil
}

func NewClient(r io.Reader, w io.Writer) (*Client, error) {
	// pr, pw := io.Pipe()
	// go io.Copy(io.MultiWriter(log.Writer(), pw), r)

	c := &Client{
		r: bufio.NewReader(r),
		w: w,
	}

	resp, err := c.tx(&Request{Cmd: reset})
	if err != nil {
		return nil, err
	}
	c.pinCount = int(resp.PinCount)

	if c.pinCount == 0 {
		return nil, errors.New("xb: no pins")
	}
	if c.pinCount > 64 {
		return nil, errors.New("xb: too many pins")
	}

	return c, nil
}

func (c *Client) tx(r *Request) (Response, error) {
	if err := WriteChunk(c.w, 'Q', r.encode()); err != nil {
		return Response{}, err
	}

	for {
		typeCode, data, err := ReadChunk(c.r)
		if err != nil {
			return Response{}, err
		}
		switch typeCode {
		case 'R':
			var resp Response
			resp.decode(data)
			if resp.Err != "" {
				return resp, errors.New(resp.Err)
			}
			return resp, nil
		case 'E':
			log.Println("xb: error:", string(data))
		case 'L':
			log.Println("xb: remote log:", string(data))
		default:
			return Response{}, fmt.Errorf("xb: unknown type code %q", typeCode)
		}
	}
}

func (c *Client) Ping() error {
	_, err := c.tx(&Request{Cmd: hello})
	return err
}

func (c *Client) PinCount() int { return c.pinCount }

func (c *Client) Pin(n int) driver.Pin {
	return &driver.PinFN{
		N:            n,
		GetFunc:      c.getPin,
		SetInputFunc: c.setInput,
		SetFunc:      c.setPin,
	}
}

func (c *Client) setInput(n int, v bool) error {
	_, err := c.tx(&Request{Cmd: setInput, Pin: uint8(n), State: v})
	return err
}

func (c *Client) setPin(n int, v bool) error {
	_, err := c.tx(&Request{Cmd: setPin, Pin: uint8(n), State: v})
	return err
}

func (c *Client) getPin(n int) (bool, error) {
	resp, err := c.tx(&Request{Cmd: getPin, Pin: uint8(n)})
	if err != nil {
		return false, err
	}
	return resp.State, nil
}

type spiClient Client

func (c *Client) SPI(cfg SPIConfig) (spi.Controller, error) {
	_, err := c.tx(&Request{Cmd: spiSetup, SPIConfig: &cfg})
	if err != nil {
		return nil, err
	}

	return (*spiClient)(c), nil
}

func (c *spiClient) SetFill(fill byte) error {
	_, err := (*Client)(c).tx(&Request{Cmd: spiSetFill, DataByte: fill})
	return err
}

func (c *spiClient) Write(data []byte) (int, error) {
	// TODO: serial buffer issue?
	if len(data) > 96 {
		n, err := c.Write(data[:96])
		if err != nil {
			return n, err
		}
		n, err = c.Write(data[96:])
		return n + 96, err
	}

	_, err := (*Client)(c).tx(&Request{Cmd: spiWrite, Data: data})
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (c *spiClient) Read(data []byte) (int, error) {
	resp, err := (*Client)(c).tx(&Request{Cmd: spiRead, ReadN: uint16(len(data))})
	if err != nil {
		return 0, err
	}
	copy(data, resp.Data)
	return len(data), nil
}

func (c *spiClient) ReadWrite(data []byte) (int, error) {
	resp, err := (*Client)(c).tx(&Request{Cmd: spiReadWrite, Data: data})
	if err != nil {
		return 0, err
	}
	copy(data, resp.Data)
	return len(data), nil
}

func (c *spiClient) ReadWriteByte(b byte) (byte, error) {
	resp, err := (*Client)(c).tx(&Request{Cmd: spiReadWriteByte, DataByte: b})
	if err != nil {
		return 0, err
	}
	return resp.DataByte, nil
}

func (c *Client) I2C(cfg I2CConfig) i2c.Bus {
	_, err := c.tx(&Request{Cmd: i2cSetup, I2CConfig: &cfg})
	if err != nil {
		return nil
	}

	return (*i2cClient)(c)
}

type i2cClient Client

func (c *i2cClient) Tx(addr uint16, w, r []byte) error {
	resp, err := (*Client)(c).tx(&Request{Cmd: i2cTx, Data: w, ReadN: uint16(len(r))})
	if err != nil {
		return err
	}
	copy(r, resp.Data)
	return nil
}
