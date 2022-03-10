package onewire

import (
	"encoding/binary"
	"errors"
	"io"
)

type Bus interface {
	Reset() bool
	WriteBit(bool)
	ReadBit() bool
}

type OneWire struct {
	Bus
}

func NewOneWire(b Bus) *OneWire {
	return &OneWire{Bus: b}
}

var (
	ErrNoDevice    = errors.New("onewire: no device found")
	ErrBadChecksum = errors.New("onewire: bad checksum")
)

func (ow *OneWire) SearchROM(alarm bool) ([]Address, error) {
	s := &searchState{alarm: alarm}
	ow.searchStart(s, 0, 0)

	return s.found, s.err
}

// ReadROM will read the 64-bit serial number of the device.
//
// Only valid when there is a single device on the bus.
func (ow *OneWire) ReadROM() (Address, error) {
	if !ow.Reset() {
		return 0, ErrNoDevice
	}

	if err := ow.WriteByte(0x33); err != nil {
		return 0, err
	}

	var a Address
	if err := binary.Read(ow, binary.BigEndian, &a); err != nil {
		return 0, err
	}

	if !a.Valid() {
		return a, ErrBadChecksum
	}

	return a, nil
}

func (ow *OneWire) Write(p []byte) (int, error) {
	if w, ok := ow.Bus.(io.Writer); ok {
		return w.Write(p)
	}

	for i, b := range p {
		if err := ow.WriteByte(b); err != nil {
			return i, err
		}
	}

	return len(p), nil
}

func (ow *OneWire) WriteByte(b byte) error {
	if bw, ok := ow.Bus.(io.ByteWriter); ok {
		return bw.WriteByte(b)
	}

	for i := 0; i < 8; i++ {
		ow.WriteBit(b&0x01 != 0)
		b >>= 1
	}

	return nil
}

func (ow *OneWire) Read(p []byte) (int, error) {
	if r, ok := ow.Bus.(io.Reader); ok {
		return r.Read(p)
	}

	var err error
	for i := range p {
		if p[i], err = ow.ReadByte(); err != nil {
			return i, err
		}
	}

	return len(p), nil
}

func (ow *OneWire) ReadByte() (byte, error) {
	if br, ok := ow.Bus.(io.ByteReader); ok {
		return br.ReadByte()
	}

	var b byte
	for i := 0; i < 8; i++ {
		b >>= 1
		if ow.ReadBit() {
			b |= 0x80
		}
	}

	return b, nil
}
