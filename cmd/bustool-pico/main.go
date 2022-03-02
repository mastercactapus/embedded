package main

import (
	"context"
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
	i2cSh := bustool.AddI2C(sh, i2cPin(machine.I2C0_SDA_PIN), i2cPin(machine.I2C0_SCL_PIN))
	bustool.AddMem(i2cSh)

	err = sh.Exec(context.Background())
	if err != nil {
		panic(err)
	}
}
