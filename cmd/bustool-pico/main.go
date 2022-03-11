package main

import (
	"machine"

	"github.com/mastercactapus/embedded/bustool"
	"github.com/mastercactapus/embedded/serial/i2c"
	"github.com/mastercactapus/embedded/serial/onewire"
	"github.com/mastercactapus/embedded/term"
)

func main() {
	err := machine.Serial.Configure(machine.UARTConfig{
		BaudRate: 115200,
	})
	if err != nil {
		panic(err)
	}

	print("\r\n")
	sh := bustool.NewShell(machine.Serial, machine.Serial)
	sh.SetNoExit(true)

	i2cSh := bustool.AddI2C(sh, i2c.I2C0())
	bustool.AddMem(i2cSh)
	bustool.AddIO(i2cSh)
	bustool.AddLCD(i2cSh)

	ow := onewire.New(onewire.NewController(19))
	owSh := sh.NewSubShell("ow", "Interact with 1-wire devices.", nil)
	owSh.AddCommand("scan", "Scan for 1-wire devices", func(ra term.RunArgs) error {
		alarm := ra.Bool(term.Flag{Name: "alarm", Short: 'a', Desc: "Filter to devices in alarm state."})
		if err := ra.Parse(); err != nil {
			return err
		}

		addrs, err := ow.SearchROM(*alarm)
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			ra.Printf("0x%x\n", addr)
		}

		return nil
	})

	err = sh.Run()
	if err != nil {
		panic(err)
	}
}
