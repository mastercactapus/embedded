package bustool

import (
	"context"
	"encoding/hex"

	"github.com/mastercactapus/embedded/i2c"
	"github.com/mastercactapus/embedded/term"
)

type ctxKey int

const (
	ctxKeyI2C ctxKey = iota
	ctxKeyMem
)

func AddI2C(sh *term.Shell, sda, scl i2c.Pin) *term.Shell {
	i2cSh := sh.NewSubShell(term.Command{Name: "i2c", Desc: "Interact with I2C devices.", Init: func(ctx context.Context, exec term.CmdFunc) error {
		if err := term.Parse(ctx).Err(); err != nil {
			return err
		}

		bus := i2c.New()
		err := bus.Configure(i2c.Config{
			SDA: sda,
			SCL: scl,
		})
		if err != nil {
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
		f := term.Parse(ctx)
		addr := f.FlagByte(term.Flag{Name: "d", Env: "DEVICE", Desc: "Device addresss.", Req: true})
		write := f.FlagBool(term.Flag{Name: "w", Desc: "Ping the write address instead."})
		if err := f.Err(); err != nil {
			return err
		}

		bus := ctx.Value(ctxKeyI2C).(*i2c.I2C)

		if write {
			return bus.PingW(addr)
		}

		return bus.Ping(addr)
	}},

	{Name: "scan", Desc: "Scan for devices.", Exec: func(ctx context.Context) error {
		if err := term.Parse(ctx).Err(); err != nil {
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
		f := term.Parse(ctx)
		addr := f.FlagByte(term.Flag{Name: "d", Env: "DEVICE", Desc: "Device addresss.", Req: true})
		reg := f.FlagByte(term.Flag{Name: "r", Def: "0", Env: "REGISTER", Desc: "Register address."})
		str := f.FlagString(term.Flag{Name: "s", Env: "DATA", Desc: "Value to write."})
		data := f.ArgByteN(term.Arg{Name: "data", Desc: "Value to write."})
		if err := f.Err(); err != nil {
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

		bus := ctx.Value(ctxKeyI2C).(*i2c.I2C)

		return bus.WriteRegister(addr, reg, data)
	}},

	{Name: "r", Desc: "Read from an I2C device register.", Exec: func(ctx context.Context) error {
		f := term.Parse(ctx)
		addr := f.FlagByte(term.Flag{Name: "d", Env: "DEVICE", Desc: "Device addresss.", Req: true})
		reg := f.FlagByte(term.Flag{Name: "r", Def: "0", Env: "REGISTER", Desc: "Register address.", Req: true})
		count := f.FlagInt(term.Flag{Name: "n", Def: "0", Desc: "Number of bytes to read."})
		if err := f.Err(); err != nil {
			return err
		}

		bus := ctx.Value(ctxKeyI2C).(*i2c.I2C)

		r := make([]byte, count)
		err := bus.ReadRegister(addr, reg, r)
		if err != nil {
			return err
		}

		term.Printer(ctx).Println(hex.Dump(r))
		return nil
	}},
}
