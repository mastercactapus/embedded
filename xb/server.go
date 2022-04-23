package xb

import (
	"bufio"
	"errors"
	"io"
	"runtime"

	"github.com/mastercactapus/embedded/driver"
	"github.com/mastercactapus/embedded/serial/i2c"
	"github.com/mastercactapus/embedded/serial/spi"
	"github.com/mastercactapus/embedded/term/ascii"
)

type Server struct {
	w io.Writer
	r *bufio.Reader

	dev  driver.Pinner
	pins []driver.Pin

	spi *spi.SoftCtrl
	i2c *i2c.I2C
}

func NewServer(r io.Reader, w io.Writer, dev driver.Pinner) *Server {
	pins := make([]driver.Pin, dev.PinCount())
	for i := range pins {
		pins[i] = dev.Pin(i)
	}
	return &Server{
		dev:  dev,
		pins: pins,
		w:    w,
		r:    bufio.NewReader(r),
	}
}

func (s *Server) logf(format string, args ...interface{}) {
	WriteChunk(s.w, 'L', []byte(ascii.Sprintf(format, args...)))
}
func (s *Server) writeResp(resp *Response) error { return WriteChunk(s.w, 'R', resp.encode()) }

func (s *Server) Serve() error {
	for {
		runtime.GC()
		typeCode, data, err := ReadChunk(s.r)
		if err != nil {
			WriteChunk(s.w, 'E', []byte(err.Error()))
			continue
		}
		var req Request
		switch typeCode {
		case 'Q':
			req.decode(data)
		default:
			WriteChunk(s.w, 'E', []byte("unknown type code"))
			continue
		}

		resp, err := s.handle(req)
		if err != nil {
			// ignore err
			s.writeResp(&Response{Err: err.Error()})
			continue
		}
		if resp == nil {
			resp = &Response{}
		}

		s.writeResp(resp)
	}
}

func (s *Server) handle(req Request) (*Response, error) {
	for {
		switch req.Cmd {
		case reset:
			return &Response{PinCount: uint8(len(s.pins))}, nil
		case hello:
			return nil, nil
		case setInput:
			return nil, s.pins[req.Pin].SetInput(req.State)
		case setPin:
			return nil, s.pins[req.Pin].Set(req.State)
		case getPin:
			state, err := s.pins[req.Pin].Get()
			return &Response{State: state}, err
		case i2cSetup:
			s.i2c = nil
			if req.I2CConfig == nil {
				return nil, errors.New("i2cSetup: missing I2CConfig")
			}
			s.i2c = i2c.New(i2c.NewSoftController(s.pins[req.I2CConfig.SDA], s.pins[req.I2CConfig.SCL]))
			return nil, nil
		case i2cTx:
			if s.i2c == nil {
				return nil, errors.New("i2cTx: i2c not initialized")
			}
			buf := make([]byte, req.ReadN)
			if err := s.i2c.Tx(req.I2CAddr, []byte(req.Data), buf); err != nil {
				return nil, err
			}
			return &Response{Data: buf}, nil
		case spiSetup:
			s.spi = nil
			if req.SPIConfig == nil {
				return nil, errors.New("spiSetup: missing SPIConfig")
			}
			cfg := &spi.Config{
				Mode: spi.Mode(req.SPIConfig.Mode),
				MOSI: s.dev.Pin(int(req.SPIConfig.MOSI)),
				MISO: s.dev.Pin(int(req.SPIConfig.MISO)),
				SCLK: s.dev.Pin(int(req.SPIConfig.SCLK)),
				Baud: int(req.SPIConfig.Baud),
			}
			spi, err := spi.NewSoftCtrl(cfg)
			if err != nil {
				return nil, err
			}

			s.spi = spi
			return nil, nil
		case spiReadWriteByte:
			if s.spi == nil {
				return nil, errors.New("spiReadWriteByte: spi not initialized")
			}

			b, err := s.spi.ReadWriteByte(req.DataByte)
			if err != nil {
				return nil, err
			}
			return &Response{DataByte: b}, nil
		case spiRead:
			if s.spi == nil {
				return nil, errors.New("spiRead: spi not initialized")
			}

			buf := make([]byte, req.ReadN)
			if _, err := s.spi.Read(buf); err != nil {
				return nil, err
			}
			return &Response{Data: buf}, nil
		case spiWrite:
			if s.spi == nil {
				return nil, errors.New("spiWrite: spi not initialized")
			}

			if _, err := s.spi.Write(req.Data); err != nil {
				return nil, err
			}
			return nil, nil
		case spiSetFill:
			if s.spi == nil {
				return nil, errors.New("spiSetFill: spi not initialized")
			}

			return nil, s.spi.SetFill(req.DataByte)
		case spiReadWrite:
			if s.spi == nil {
				return nil, errors.New("spiReadWrite: spi not initialized")
			}

			if _, err := s.spi.ReadWrite(req.Data); err != nil {
				return nil, err
			}
			return &Response{Data: req.Data}, nil
		}
	}
}
