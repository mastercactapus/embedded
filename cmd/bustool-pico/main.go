package main

import (
	"machine"

	"github.com/mastercactapus/embedded/bustool"
)

func main() {
	err := machine.Serial.Configure(machine.UARTConfig{
		BaudRate: 115200,
	})
	if err != nil {
		panic(err)
	}

	sh := bustool.NewShell(&fixReader{machine.Serial}, machine.Serial)
	bustool.AddI2C(sh, i2cPin(machine.I2C0_SDA_PIN), i2cPin(machine.I2C0_SCL_PIN))
	machine.LED.High()
	err = sh.Exec()
	if err != nil {
		panic(err)
	}
}
