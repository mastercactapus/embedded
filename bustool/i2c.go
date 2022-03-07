package bustool

import (
	"encoding/hex"

	"github.com/mastercactapus/embedded/bus/i2c"
	"github.com/mastercactapus/embedded/term"
)

func AddI2C(sh *term.Shell2, bus *i2c.I2C) *term.Shell2 {
	i2cSh := sh.NewSubShell(term.Command2{Name: "i2c", Desc: "Interact with I2C devices.", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		r.Set("i2c", bus)
		return nil
	}})

	for _, c := range i2cCommands {
		i2cSh.AddCommand(c)
	}
	return i2cSh
}

var i2cCommands = []term.Command2{
	{Name: "ping", Desc: "Ping a device.", Exec: func(r term.RunArgs) error {
		addr := r.Byte(term.Flag{Name: "dev", Short: 'd', Env: "DEV", Desc: "Device addresss.", Req: true})
		write := r.Bool(term.Flag{Short: 'w', Desc: "Ping the write address instead."})
		if err := r.Parse(); err != nil {
			return err
		}

		bus := r.Get("i2c").(*i2c.I2C)

		if *write {
			return bus.PingW(*addr)
		}

		return bus.Ping(*addr)
	}},

	{Name: "scan", Desc: "Scan for devices.", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		bus := r.Get("i2c").(*i2c.I2C)

		for i := 0; i < 127; i++ {
			canRead := bus.Ping(byte(i)) == nil
			canWrite := bus.PingW(byte(i)) == nil
			if !canRead && !canWrite {
				continue
			}

			r.Printf("0x%02x: ", i)
			switch {
			case canRead && canWrite:
				r.Printf("RW")
			case canRead && !canWrite:
				r.Printf("RO")
			case !canRead && canWrite:
				r.Printf("WO")
			}

			id, err := bus.DeviceID(byte(i))
			if err == nil {
				r.Printf(" ID=%x", id)
			}
			r.Println()
		}

		return nil
	}},

	{Name: "w", Desc: "Write to an I2C device register.", Exec: func(r term.RunArgs) error {
		addr := r.Byte(term.Flag{Name: "dev", Short: 'd', Env: "DEV", Desc: "Device addresss.", Req: true})
		reg := r.Byte(term.Flag{Name: "reg", Short: 'r', Def: "0", Env: "REG", Desc: "Register address."})
		data := r.Bytes(term.Flag{Name: "data", Short: 'b', Desc: "Write bytes (comma separated).", Req: true})
		if err := r.Parse(); err != nil {
			return err
		}

		bus := r.Get("i2c").(*i2c.I2C)
		return bus.WriteRegister(*addr, *reg, *data)
	}},

	{Name: "r", Desc: "Read from an I2C device register.", Exec: func(r term.RunArgs) error {
		addr := r.Byte(term.Flag{Name: "dev", Short: 'd', Env: "DEV", Desc: "Device addresss.", Req: true})
		reg := r.Byte(term.Flag{Name: "reg", Short: 'r', Def: "0", Env: "REG", Desc: "Register address."})
		count := r.Int(term.Flag{Name: "n", Def: "0", Desc: "Number of bytes to read."})
		if err := r.Parse(); err != nil {
			return err
		}

		bus := r.Get("i2c").(*i2c.I2C)

		data := make([]byte, *count)
		err := bus.ReadRegister(*addr, *reg, data)
		if err != nil {
			return err
		}

		r.Println(hex.Dump(data))
		return nil
	}},

	{Name: "tx", Desc: "Read/write to an I2C device.", Exec: func(r term.RunArgs) error {
		addr := r.Byte(term.Flag{Name: "dev", Short: 'd', Env: "DEV", Desc: "Device addresss.", Req: true})
		count := r.Int(term.Flag{Name: "n", Def: "0", Desc: "Number of bytes to read."})
		data := r.Bytes(term.Flag{Name: "data", Short: 'b', Desc: "Write bytes (comma separated)."})
		if err := r.Parse(); err != nil {
			return err
		}

		bus := r.Get("i2c").(*i2c.I2C)

		var rData []byte
		if *count > 0 {
			rData = make([]byte, *count)
		}
		err := bus.Tx(uint16(*addr), *data, rData)
		if err != nil {
			return err
		}

		if *count > 0 {
			r.Println(hex.Dump(rData))
		}
		return nil
	}},
}
