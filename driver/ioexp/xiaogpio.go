package ioexp

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/mastercactapus/embedded/driver"
	"github.com/mastercactapus/embedded/serial/i2c"
)

type XIAO struct {
	rw io.ReadWriter

	InputPins   *Register16
	OutputState *Register16
	ReadState   *Register16
}

func NewXIAO(rw io.ReadWriter) *XIAO {
	return &XIAO{
		rw:          rw,
		InputPins:   NewRegister16(rw, 'D'),
		OutputState: NewRegister16(rw, 'V'),
		ReadState:   NewRegister16(rw, 'R'),
	}
}

func (XIAO) PinCount() int { return 11 }

func (x *XIAO) Flush() error {
	if err := x.InputPins.Flush(); err != nil {
		return err
	}
	if err := x.OutputState.Flush(); err != nil {
		return err
	}

	return nil
}

// Ping will test that the device is responding.
func (x *XIAO) Ping() error {
	_, err := x.rw.Write([]byte{'?'})
	if err != nil {
		return err
	}

	var buf [2]byte
	_, err = io.ReadFull(x.rw, buf[:])
	if err != nil {
		return err
	}

	if string(buf[:]) != "OK" {
		return errors.New("XIAO: ping failed")
	}

	return nil
}

func (x *XIAO) Pin(n int) driver.Pin {
	return &driver.PinFN{
		N:            n,
		GetFunc:      x.ReadState.Get,
		SetInputFunc: x.InputPins.Set,
		SetFunc:      x.OutputState.Set,
	}
}

func (x *XIAO) BufferedPin(n int) driver.Pin {
	return &driver.PinFN{
		N:            n,
		GetFunc:      x.ReadState.GetBuf,
		SetInputFunc: x.InputPins.SetBuf,
		SetFunc:      x.OutputState.SetBuf,
	}
}

func (x *XIAO) I2C(sda, scl int) (i2c.Bus, error) {
	var args struct {
		Cmd      byte
		SDA, SCL byte
	}
	args.Cmd = 'I'
	args.SDA = byte(sda)
	args.SCL = byte(scl)
	err := binary.Write(x.rw, binary.LittleEndian, args)
	if err != nil {
		return nil, err
	}

	return (*xiaoI2C)(x), nil
}

type xiaoI2C XIAO

func (x *xiaoI2C) Tx(addr uint16, w, r []byte) error {
	var args struct {
		Cmd    byte
		Addr   uint16
		ReadN  uint8
		WriteN uint8
	}
	args.Cmd = 'i'
	args.Addr = addr
	args.ReadN = uint8(len(r))
	args.WriteN = uint8(len(w))
	err := binary.Write(x.rw, binary.LittleEndian, args)
	if err != nil {
		return err
	}
	_, err = x.rw.Write(w)
	if err != nil {
		return err
	}

	_, err = io.ReadFull(x.rw, r)
	if err != nil {
		return err
	}

	return nil
}
