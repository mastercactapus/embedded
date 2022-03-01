package bustool

import (
	"encoding/hex"
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

	for _, c := range i2cCommands {
		i2cSh.AddCommand(c)
	}
}

var i2cCommands = []term.Command{
	{Name: "ping", Desc: "Ping a device.", Exec: func(c *term.CmdCtx) error {
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
	}},

	{Name: "scan", Desc: "Scan for devices.", Exec: func(c *term.CmdCtx) error {
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

			p.Printf("0x%02x: ", i)
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
	}},

	{Name: "w", Desc: "Write to an I2C device register.", Exec: func(c *term.CmdCtx) error {
		addr := c.FlagByte(term.Flag{Name: "d", Env: "DEVICE", Desc: "Device addresss.", Req: true})
		reg := c.FlagByte(term.Flag{Name: "r", Def: "0", Env: "REGISTER", Desc: "Register address."})
		str := c.FlagString(term.Flag{Name: "s", Env: "DATA", Desc: "Value to write."})
		data := c.ArgByteN(term.Arg{Name: "data", Desc: "Value to write."})
		if err := c.Parse(); err != nil {
			return err
		}

		if len(str) > 0 && len(data) > 0 {
			return term.UsageError("cannot specify both -s and data")
		}
		if len(str) > 0 {
			data = []byte(str)
		}
		if len(data) == 0 {
			return term.UsageError("data or -s required")
		}

		bus, ok := c.Get("i2c").(*i2c.I2C)
		if !ok {
			return fmt.Errorf("i2c: not available")
		}

		_, err := bus.WriteRegister(addr, reg, data)
		return err
	}},

	{Name: "r", Desc: "Read from an I2C device register.", Exec: func(c *term.CmdCtx) error {
		addr := c.FlagByte(term.Flag{Name: "d", Env: "DEVICE", Desc: "Device addresss.", Req: true})
		reg := c.FlagByte(term.Flag{Name: "r", Def: "0", Env: "REGISTER", Desc: "Register address.", Req: true})
		count := c.FlagInt(term.Flag{Name: "n", Def: "0", Desc: "Number of bytes to read."})
		if err := c.Parse(); err != nil {
			return err
		}

		bus, ok := c.Get("i2c").(*i2c.I2C)
		if !ok {
			return fmt.Errorf("i2c: not available")
		}

		r := make([]byte, count)
		_, err := bus.ReadRegister(addr, reg, r)
		if err != nil {
			return err
		}

		c.Printer().Println(hex.Dump(r))
		return nil
	}},
}
