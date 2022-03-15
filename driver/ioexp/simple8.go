package ioexp

import (
	"io"

	"github.com/mastercactapus/embedded/driver"
)

type Simple8 struct {
	rw io.ReadWriter

	State    uint8
	readData uint8
}

func NewSimple8(rw io.ReadWriter) *Simple8 {
	return &Simple8{rw: rw}
}

func (Simple8) PinCount() int { return 8 }

func (s *Simple8) Pin(n int) driver.Pin {
	return &driver.PinFN{
		N:            n,
		SetInputFunc: s.setPin,
		SetFunc:      s.setPin,
		GetFunc:      s.getPin,
	}
}

func (s *Simple8) BufferedPin(n int) driver.Pin {
	return &driver.PinFN{
		N:            n,
		SetInputFunc: s.setPinB,
		SetFunc:      s.setPinB,
		GetFunc:      s.getPinB,
	}
}

func (p *Simple8) getPinB(n int) (bool, error) {
	return p.readData&(1<<uint8(n)) != 0, nil
}

func (p *Simple8) getPin(n int) (v bool, err error) {
	err = p.Refresh()
	if err != nil {
		return false, err
	}
	return p.getPinB(n)
}

func (s *Simple8) setPinB(n int, v bool) error {
	if v {
		s.State |= 1 << uint8(n)
	} else {
		s.State &= ^(1 << uint8(n))
	}
	return nil
}

func (s *Simple8) setPin(n int, v bool) error {
	s.setPinB(n, v)
	return s.Flush()
}

func (s *Simple8) Refresh() (err error) {
	s.readData, err = s.read()
	return err
}

func (s *Simple8) Flush() error {
	if bw, ok := s.rw.(io.ByteWriter); ok {
		return bw.WriteByte(s.State)
	}

	_, err := s.rw.Write([]byte{s.State})
	return err
}

func (p *Simple8) read() (b uint8, err error) {
	if br, ok := p.rw.(io.ByteReader); ok {
		return br.ReadByte()
	}

	buf := make([]byte, 1)
	_, err = p.rw.Read(buf[:])
	return buf[0], err
}
