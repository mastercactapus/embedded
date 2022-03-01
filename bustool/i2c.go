package bustool

import (
	"fmt"

	"github.com/mastercactapus/embedded/i2c"
	"github.com/mastercactapus/embedded/term"
)

func AddI2C(sh *term.Shell, sda, scl i2c.Pin) {
	i2cSh := sh.NewSubShell(term.Command{Name: "i2c", Desc: "Interact with I2C devices.", Exec: func(c *term.CmdCtx) error {
		bus := i2c.New()
		bus.Configure(i2c.Config{
			SDA: sda,
			SCL: scl,
		})
		c.Set("i2c", bus)
		return nil
	}})

	i2cSh.AddCommand(term.Command{Name: "ping", Desc: "Ping a device.", Exec: func(c *term.CmdCtx) error {
		addr := c.FlagByte(term.Flag{Name: "d", Env: "DEVICE", Desc: "Device addresss.", Req: true})
		write := c.FlagBool(term.Flag{Name: "w", Desc: "Ping the write address instead."})
		if err := c.Parse(); err != nil {
			return err
		}

		bus, ok := c.Get("i2c").(*i2c.I2C)
		if !ok {
			return fmt.Errorf("i2c: not available")
		}

		if write {
			return bus.PingW(addr)
		}

		return bus.Ping(addr)
	}})

	i2cSh.AddCommand(term.Command{Name: "scan", Desc: "Scan for devices.", Exec: func(c *term.CmdCtx) error {
		if err := c.Parse(); err != nil {
			return err
		}

		bus, ok := c.Get("i2c").(*i2c.I2C)
		if !ok {
			return fmt.Errorf("i2c: not available")
		}

		p := c.Printer()
		for i := 0; i < 127; i++ {
			canRead := bus.Ping(byte(i)) == nil
			canWrite := bus.PingW(byte(i)) == nil
			if !canRead && !canWrite {
				continue
			}

			p.Printf("%02x: ", i)
			switch {
			case canRead && canWrite:
				p.Printf("RW")
			case canRead && !canWrite:
				p.Printf("RO")
			case !canRead && canWrite:
				p.Printf("WO")
			}

			id, err := bus.DeviceID(byte(i))
			if err == nil {
				p.Printf(" ID=%x", id)
			}
			p.Println()
		}

		return nil
	}})
}
