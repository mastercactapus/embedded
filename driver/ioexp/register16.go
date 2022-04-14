package ioexp

import (
	"io"

	"github.com/mastercactapus/embedded/serial"
)

type Register16 struct {
	rw   io.ReadWriter
	addr uint8

	State uint16

	invert uint16
}

func NewRegister16(rw io.ReadWriter, addr uint8) *Register16 {
	return &Register16{addr: addr, rw: rw}
}

func (r *Register16) SetInvert(n int, v bool) {
	if v {
		r.invert |= 1 << n
	} else {
		r.invert &^= 1 << n
	}
}

func (r *Register16) Flush() error {
	_, err := r.rw.Write([]byte{r.addr, byte(r.State ^ r.invert), byte((r.State ^ r.invert) >> 8)})
	return err
}

func (r *Register16) SetBuf(n int, v bool) error {
	if v {
		r.State |= 1 << n
	} else {
		r.State &^= 1 << n
	}
	return nil
}

func (r *Register16) Set(n int, v bool) error {
	if err := r.SetBuf(n, v); err != nil {
		return err
	}
	return r.Flush()
}

func (r *Register16) GetBuf(n int) (bool, error) {
	return r.State&(1<<n) != 0, nil
}

func (r *Register16) Get(n int) (bool, error) {
	if err := r.Refresh(); err != nil {
		return false, err
	}
	return r.GetBuf(n)
}

func (r *Register16) Refresh() error {
	var buf [2]byte
	err := serial.Tx(r.rw, []byte{r.addr}, buf[:])
	if err != nil {
		return err
	}
	r.State = (uint16(buf[0]) | uint16(buf[1])<<8) ^ r.invert
	return nil
}
