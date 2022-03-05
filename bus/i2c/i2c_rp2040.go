//go:build pico
// +build pico

package i2c

import (
	"device"
	"machine"
)

func wait() {
	for i := 0; i < 50; i++ {
		device.Asm("nop")
	}
}

type i2cPin machine.Pin

var _ Pin = i2cPin(0)

func (p i2cPin) PullupHigh() {
	machine.Pin(p).Configure(machine.PinConfig{Mode: machine.PinInputPullup})
}

func (p i2cPin) OutputLow() {
	machine.Pin(p).Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.Pin(p).Low()
}

func (p i2cPin) Get() bool {
	return machine.Pin(p).Get()
}

func I2C0() (*I2C, error) {
	n := New()
	err := n.Configure(Config{
		SDA: i2cPin(machine.I2C0_SDA_PIN),
		SCL: i2cPin(machine.I2C0_SCL_PIN),
	})
	if err != nil {
		return nil, err
	}
	return n, nil
}
