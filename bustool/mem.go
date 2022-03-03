package bustool

import (
	"context"
	"encoding/hex"

	"github.com/mastercactapus/embedded/driver/eeprom"
	"github.com/mastercactapus/embedded/term"
)

func AddMem(sh *term.Shell) *term.Shell {
	memSh := sh.NewSubShell(term.Command{Name: "mem", Desc: "Interact with an AT24Cxx-compatible EEPROM device over I2C.", Init: func(ctx context.Context, exec term.CmdFunc) error {
		f := term.Parse(ctx)
		addr := f.FlagByte(term.Flag{Name: "d", Def: "0x50", Env: "DEVICE", Desc: "Device addresss.", Req: true})
		mem := f.FlagInt(term.Flag{Name: "m", Def: "128", Desc: "Memory size in bits.", Req: true})
		if f.Err() != nil {
			return f.Err()
		}

		var dev *eeprom.Device
		bus := ctx.Value(ctxKeyI2C).(eeprom.I2C)

		switch mem {
		case 128:
			dev = eeprom.NewAT24C01(bus, addr)
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
		f := term.Parse(ctx)
		start := f.FlagInt(term.Flag{Name: "p", Def: "0", Desc: "Position to start from.", Req: true})
		count := f.FlagInt(term.Flag{Name: "n", Def: "0", Desc: "Number of bytes to read, if zero read to end."})
		if err := f.Err(); err != nil {
			return err
		}

		mem := ctx.Value(ctxKeyMem).(*eeprom.Device)
		if count == 0 {
			count = mem.Size() - start
		}

		data := make([]byte, count)
		_, err := mem.ReadAt(data, int64(start))
		if err != nil {
			return err
		}

		term.Printer(ctx).Print(hex.Dump(data))
		return nil
	}},
	{Name: "w", Desc: "Write device data.", Exec: func(ctx context.Context) error {
		f := term.Parse(ctx)
		start := f.FlagInt(term.Flag{Name: "p", Def: "0", Desc: "Position to start from.", Req: true})
		str := f.FlagString(term.Flag{Name: "s", Env: "DATA", Desc: "Value to write."})
		data := f.ArgByteN(term.Arg{Name: "data", Desc: "Value to write."})
		if err := f.Err(); err != nil {
			return err
		}

		if len(str) > 0 && len(data) > 0 {
			return f.UsageError("cannot specify both -s and data")
		}
		if len(str) > 0 {
			data = []byte(str)
		}
		if len(data) == 0 {
			return f.UsageError("data or -s required")
		}

		mem := ctx.Value(ctxKeyMem).(*eeprom.Device)

		_, err := mem.WriteAt(data, int64(start))
		if err != nil {
			println(err.Error())
			return err
		}

		return nil
	}},
}
