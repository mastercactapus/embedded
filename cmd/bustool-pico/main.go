package main

import (
	"machine"

	"github.com/mastercactapus/embedded/bus/i2c"
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
	sh.SetNoExit(true)
	bus, err := i2c.I2C0()
	if err != nil {
		panic(err)
	}
	i2cSh := bustool.AddI2C(sh, bus)
	bustool.AddMem(i2cSh)
	bustool.AddIO(i2cSh)
	bustool.AddLCD(i2cSh)

	err = sh.Run()
	if err != nil {
		panic(err)
	}
}
