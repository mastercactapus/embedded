package bustool

import (
	"context"
	"encoding/hex"
	"math/rand"

	"github.com/mastercactapus/embedded/driver/eeprom"
	"github.com/mastercactapus/embedded/term"
)

func AddMem(sh *term.Shell) *term.Shell {
	memSh := sh.NewSubShell(term.Command{Name: "mem", Desc: "Interact with an AT24Cxx-compatible EEPROM device over I2C.", Init: func(ctx context.Context, exec term.CmdFunc) error {
		f := term.Flags(ctx)
		addr := f.Byte(term.Flag{Name: "dev", Short: 'd', Def: "0x50", Env: "DEV", Desc: "Device addresss.", Req: true})
		mem := f.Int(term.Flag{Name: "size", Short: 'm', Def: "1", Desc: "Memory size in kbits.", Req: true})
		if err := f.Parse(); err != nil {
			return err
		}

		var dev *eeprom.Device
		bus := ctx.Value(ctxKeyI2C).(eeprom.I2C)

		switch *mem {
		case 1:
			dev = eeprom.NewAT24C01(bus, *addr)
		default:

			return f.UsageError("unsupported memory size %d", mem)
		}

		return exec(context.WithValue(ctx, ctxKeyMem, dev))
	}})

	for _, c := range memCommands {
		memSh.AddCommand(c)
	}
	return memSh
}

var memCommands = []term.Command{
	{Name: "r", Desc: "Read device data.", Exec: func(ctx context.Context) error {
		f := term.Flags(ctx)
		start := f.Int(term.Flag{Name: "p", Def: "0", Desc: "Position to start from.", Req: true})
		count := f.Int(term.Flag{Name: "n", Def: "0", Desc: "Number of bytes to read, if zero read to end."})
		if err := f.Parse(); err != nil {
			return err
		}

		mem := ctx.Value(ctxKeyMem).(*eeprom.Device)
		if *count == 0 {
			*count = mem.Size() - *start
		}
		if *count <= 0 {
			return nil
		}

		data := make([]byte, *count)
		_, err := mem.ReadAt(data, int64(*start))
		if err != nil {
			return err
		}

		term.Printer(ctx).Print(hex.Dump(data))
		return nil
	}},
	{Name: "w", Desc: "Write device data.", Exec: func(ctx context.Context) error {
		f := term.Flags(ctx)
		start := f.Int(term.Flag{Name: "p", Def: "0", Desc: "Position to start from.", Req: true})
		data := f.BinaryArgs(term.Arg{Name: "data", Desc: "Value to write."})
		if err := f.Parse(); err != nil {
			return err
		}

		mem := ctx.Value(ctxKeyMem).(*eeprom.Device)

		_, err := mem.WriteAt(*data, int64(*start))
		if err != nil {
			return err
		}

		return nil
	}},
	{Name: "format", Desc: "Clear all data.", Exec: func(ctx context.Context) error {
		f := term.Flags(ctx)
		start := f.Int(term.Flag{Name: "p", Def: "0", Desc: "Position to start from.", Req: true})
		count := f.Int(term.Flag{Name: "n", Def: "0", Desc: "Number of bytes to wipe, if zero clear to end."})
		value := f.Byte(term.Flag{Name: "v", Def: "0xff", Desc: "Value to write."})
		rnd := f.Bool(term.Flag{Name: "random", Desc: "Fill with random data."})
		if err := f.Parse(); err != nil {
			return err
		}

		mem := ctx.Value(ctxKeyMem).(*eeprom.Device)

		if *count == 0 {
			*count = mem.Size() - *start
		}
		if *count <= 0 {
			return nil
		}

		data := make([]byte, *count)
		for i := range data {
			if *rnd {
				data[i] = byte(rand.Intn(256))
			} else {
				data[i] = *value
			}
		}

		_, err := mem.WriteAt(data, int64(*start))
		if err != nil {
			return err
		}

		return nil
	}},
}
