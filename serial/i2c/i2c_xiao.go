//go:build xiao
// +build xiao

package i2c

import (
	"machine"

	"github.com/mastercactapus/embedded/driver"
)

func I2C0() *I2C {
	return New(NewSoftController(
		driver.FromMachine(machine.SDA_PIN),
		driver.FromMachine(machine.SCL_PIN),
	))
}
