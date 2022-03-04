package ioexp

import (
	"time"
)

type I2C interface {
	Tx(addr uint16, w, r []byte) error
}

func NewPCF8574(i2c I2C, addr uint8) *Device {
	return &Device{i2c: i2c, addr: addr}
}

type Device struct {
	i2c  I2C
	addr uint8
}

type Pins uint8

// Pins returns the current pin state.
func (d *Device) Pins() (Pins, error) {
	var buf [1]byte
	err := d.i2c.Tx(uint16(d.addr), nil, buf[:])
	time.Sleep(5 * time.Microsecond)
	if err != nil {
		return 0, err
	}

	return Pins(buf[0]), nil
}

// SetPins writes the new pin state.
func (d *Device) SetPins(p Pins) error {
	var buf [1]byte
	buf[0] = byte(p)
	err := d.i2c.Tx(uint16(d.addr), buf[:], nil)
	time.Sleep(5 * time.Microsecond)
	if err != nil {
		return err
	}

	return nil
}

// Get returns true if the pin is HIGH, false if LOW.
func (p Pins) Get(n int) bool {
	return (p & (1 << n)) != 0
}

// Sets the numbered pin to HIGH/Input for true or LOW/Output for false.
func (p *Pins) Set(n int, v bool) {
	if v {
		*p |= (1 << n)
	} else {
		*p &= ^(1 << n)
	}
}
