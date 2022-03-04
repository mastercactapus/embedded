package ioexp

import (
	"github.com/mastercactapus/embedded/i2c"
)

// NewPCF8574 is a convenience method that returns a PinReadWriter for a PCF8574-compatible I2C device.
func NewPCF8574(bus i2c.Bus, addr uint16) PinReadWriter {
	return NewSimple8Bit(i2c.NewDevice(bus, addr))
}
