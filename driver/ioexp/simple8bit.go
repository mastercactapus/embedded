package ioexp

import (
	"io"
)

// Simple8Bit is an I/O expander with 8 pins without registers.
type Simple8Bit struct {
	rw io.ReadWriter
	m  *PinMasker
}

var (
	_ PinReadWriter = (*Simple8Bit)(nil)
	_ InputSetter   = (*Simple8Bit)(nil)
)

// NewSimple8Bit returns a device that simply reads and writes byte values
// to control the pins.
func NewSimple8Bit(rw io.ReadWriter) *Simple8Bit {
	return &Simple8Bit{rw: rw, m: NewPinMasker(8)}
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

// WritePins writes the given pin state to the device.
//
// High values also changes the pins to input mode.
// Low values also changes the pins to output mode.
func (dev *Simple8Bit) WritePins(pins Valuer) error {
	dev.m.Set(pins)
	b := PinByte(pins)

	if bw, ok := dev.rw.(io.ByteWriter); ok {
		return bw.WriteByte(b)
	}

	_, err := dev.rw.Write([]byte{byte(b)})
	return err
}

func (dev *Simple8Bit) WritePinsMask(pins, mask Valuer) error {
	return dev.m.ApplyFn(pins, mask, dev.WritePins)
}

// SetInputPins is identical to WritePins as the device
// only has HIGH/INPUT and LOW/OUTPUT states for each pin.
func (dev *Simple8Bit) SetInputPins(pins Valuer) error {
	return dev.WritePins(pins)
}

func (dev *Simple8Bit) SetInputPinsMask(pins, mask Valuer) error {
	return dev.WritePinsMask(pins, mask)
}

// TODO: add software invert?
