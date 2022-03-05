package ioexp

import (
	"github.com/mastercactapus/embedded/i2c"
)

// NewPCF8574 is a convenience method that returns a PinReadWriter for a PCF8574-compatible I2C device.
//
// Default address is 0x20.
func NewPCF8574(bus i2c.Bus, addr uint16) PinReadWriter {
	if addr == 0 {
		addr = 0x20
	}
	return NewSimple8Bit(i2c.NewDevice(bus, addr))
}
