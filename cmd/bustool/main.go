package main

import (
	"time"

	"github.com/mastercactapus/embedded/bustool"
	"github.com/mastercactapus/embedded/driver/ioexp"
	"github.com/mastercactapus/embedded/serial/i2c"
	"github.com/mastercactapus/embedded/serial/onewire"
	"github.com/mastercactapus/embedded/term"
)

func main() {
	sh := bustool.NewShell(configIO())
	sh.SetNoExit(true)

	sh.AddCommand("test", "", func(ra term.RunArgs) error {
		baseBus := i2c.I2C0()
		mcp := ioexp.NewPCF8574(baseBus, 0x20)

		sr := ioexp.NewSN74HC595(ioexp.SN74HC595Config{
			SER:   mcp.Pin(0),
			RCLK:  mcp.Pin(1),
			SRCLK: mcp.Pin(2),
		})
		err := sr.Configure(0)
		if err != nil {
			return err
		}

		for {
			for i := 0; i < 8; i++ {
				if !ra.WaitForInterrupt() {
					return nil
				}
				ra.Println("Write", i)
				sr.Pin(i).High()
				time.Sleep(1 * time.Second)
				sr.Pin(i).Low()
			}
		}
	})

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

	panic(sh.Run())
}
