package bustool

import (
	"context"
	"encoding/hex"

	"github.com/mastercactapus/embedded/bus/i2c"
	"github.com/mastercactapus/embedded/term"
)

type ctxKey int

const (
	ctxKeyI2C ctxKey = iota
	ctxKeyMem
	ctxKeyIO
	ctxKeyLCD
)

func AddI2C(sh *term.Shell, bus *i2c.I2C) *term.Shell {
	i2cSh := sh.NewSubShell(term.Command{Name: "i2c", Desc: "Interact with I2C devices.", Init: func(ctx context.Context, exec term.CmdFunc) error {
		if err := term.Flags(ctx).Parse(); err != nil {
			return err
		}

		return exec(context.WithValue(ctx, ctxKeyI2C, bus))
	}})

	for _, c := range i2cCommands {
		i2cSh.AddCommand(c)
	}
	return i2cSh
}

var i2cCommands = []term.Command{
	{Name: "ping", Desc: "Ping a device.", Exec: func(ctx context.Context) error {
		f := term.Flags(ctx)
		addr := f.Byte(term.Flag{Name: "dev", Short: 'd', Env: "DEV", Desc: "Device addresss.", Req: true})
		write := f.Bool(term.Flag{Short: 'w', Desc: "Ping the write address instead."})
		if err := f.Parse(); err != nil {
			return err
		}

		bus := ctx.Value(ctxKeyI2C).(*i2c.I2C)

		if *write {
			return bus.PingW(*addr)
		}

		return bus.Ping(*addr)
	}},

	{Name: "scan", Desc: "Scan for devices.", Exec: func(ctx context.Context) error {
		if err := term.Flags(ctx).Parse(); err != nil {
			return err
		}

		bus := ctx.Value(ctxKeyI2C).(*i2c.I2C)

		p := term.Printer(ctx)
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

	{Name: "w", Desc: "Write to an I2C device register.", Exec: func(ctx context.Context) error {
		f := term.Flags(ctx)
		addr := f.Byte(term.Flag{Name: "dev", Short: 'd', Env: "DEV", Desc: "Device addresss.", Req: true})
		reg := f.Byte(term.Flag{Name: "reg", Short: 'r', Def: "0", Env: "REG", Desc: "Register address."})
		data := f.Bytes(term.Flag{Name: "data", Short: 'b', Desc: "Write bytes (comma separated).", Req: true})
		if err := f.Parse(); err != nil {
			return err
		}

		bus := ctx.Value(ctxKeyI2C).(*i2c.I2C)
		return bus.WriteRegister(*addr, *reg, *data)
	}},

	{Name: "r", Desc: "Read from an I2C device register.", Exec: func(ctx context.Context) error {
		f := term.Flags(ctx)
		addr := f.Byte(term.Flag{Name: "dev", Short: 'd', Env: "DEV", Desc: "Device addresss.", Req: true})
		reg := f.Byte(term.Flag{Name: "reg", Short: 'r', Def: "0", Env: "REG", Desc: "Register address."})
		count := f.Int(term.Flag{Name: "n", Def: "0", Desc: "Number of bytes to read."})
		if err := f.Parse(); err != nil {
			return err
		}

		bus := ctx.Value(ctxKeyI2C).(*i2c.I2C)

		r := make([]byte, *count)
		err := bus.ReadRegister(*addr, *reg, r)
		if err != nil {
			return err
		}

		term.Printer(ctx).Println(hex.Dump(r))
		return nil
	}},

	{Name: "tx", Desc: "Read/write to an I2C device.", Exec: func(ctx context.Context) error {
		f := term.Flags(ctx)
		addr := f.Byte(term.Flag{Name: "dev", Short: 'd', Env: "DEV", Desc: "Device addresss.", Req: true})
		count := f.Int(term.Flag{Name: "n", Def: "0", Desc: "Number of bytes to read."})
		data := f.Bytes(term.Flag{Name: "data", Short: 'b', Desc: "Write bytes (comma separated)."})
		if err := f.Parse(); err != nil {
			return err
		}

		bus := ctx.Value(ctxKeyI2C).(*i2c.I2C)

		var r []byte
		if *count > 0 {
			r = make([]byte, *count)
		}
		err := bus.Tx(uint16(*addr), *data, r)
		if err != nil {
			return err
		}

		if *count > 0 {
			term.Printer(ctx).Println(hex.Dump(r))
		}
		return nil
	}},
}
