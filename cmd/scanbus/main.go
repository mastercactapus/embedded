package main

import (
	"machine"

	"github.com/mastercactapus/embedded/driver"
	"github.com/mastercactapus/embedded/serial/i2c"
	"github.com/mastercactapus/embedded/term"
)

func main() {
	machine.Serial.Configure(machine.UARTConfig{
		BaudRate: 115200,
	})
	sh := term.NewRootShell("scanbus", "Scan for bus devices.", machine.Serial, machine.Serial)
	sh.SetNoExit(true)
	sh.AddCommand("scani2c", "Scan for i2c devices.", func(ra term.RunArgs) error {
		sda := ra.Int(term.Flag{Name: "sda", Short: 'd', Desc: "SDA pin", Req: true})
		scl := ra.Int(term.Flag{Name: "scl", Short: 'c', Desc: "SCL pin", Req: true})
		if err := ra.Parse(); err != nil {
			return err
		}

		ctrl := i2c.NewSoftController(driver.FromMachine(machine.Pin(*sda)), driver.FromMachine(machine.Pin(*scl)))

		bus := i2c.New(ctrl)

		for i := 0; i < 127; i++ {
			canRead := bus.Ping(byte(i)) == nil
			canWrite := bus.PingW(byte(i)) == nil
			if !canRead && !canWrite {
				continue
			}

			ra.Printf("0x%02x: ", i)
			switch {
			case canRead && canWrite:
				ra.Printf("RW")
			case canRead && !canWrite:
				ra.Printf("RO")
			case !canRead && canWrite:
				ra.Printf("WO")
			}

			id, err := bus.DeviceID(byte(i))
			if err == nil {
				ra.Printf(" ID=%x", id)
			}
			ra.Println()
		}
		return nil
	})

	panic(sh.Run())
}
