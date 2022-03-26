package ioexp

import (
	"errors"
	"io"

	"github.com/mastercactapus/embedded/driver"
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
