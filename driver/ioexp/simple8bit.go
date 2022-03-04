package ioexp

import "io"

// Simple8Bit is a PCF8574-compatible I/O expander with 8 pins.
type Simple8Bit struct {
	rw io.ReadWriter
}

// NewSimple8Bit returns a device that simply reads and writes byte values
// to control the pins.
func NewSimple8Bit(rw io.ReadWriter) *Simple8Bit {
	return &Simple8Bit{rw: rw}
}

func (Simple8Bit) PinCount() int { return 8 }

func (dev *Simple8Bit) ReadPins() (PinState, error) {
	if br, ok := dev.rw.(io.ByteReader); ok {
		b, err := br.ReadByte()
		if err != nil {
			return nil, err
		}
		return (*Pin8)(&b), nil
	}

	var buf [1]byte
	_, err := dev.rw.Read(buf[:])
	if err != nil {
		return nil, err
	}
	return (*Pin8)(&buf[0]), nil
}

func (dev *Simple8Bit) WritePins(pins PinState) error {
	var b byte
	switch t := pins.(type) {
	case *Pin8:
		b = byte(*t)
	default:
		var p Pin8
		CopyState(&p, pins)
		b = byte(p)
	}

	if bw, ok := dev.rw.(io.ByteWriter); ok {
		return bw.WriteByte(b)
	}

	_, err := dev.rw.Write([]byte{byte(b)})
	return err
}
