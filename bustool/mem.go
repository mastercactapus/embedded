package bustool

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/mastercactapus/embedded/driver/eeprom"
	"github.com/mastercactapus/embedded/i2c"
	"github.com/mastercactapus/embedded/term"
)

func AddMem(sh *term.Shell) *term.Shell {
	memSh := sh.NewSubShell(term.Command{Name: "mem", Desc: "Interact with an AT24Cxx-compatible EEPROM device over I2C.", Init: func(ctx context.Context, exec term.CmdFunc) error {
		f := term.Parse(ctx)
		addr := f.FlagByte(term.Flag{Name: "d", Env: "DEVICE", Desc: "Device addresss.", Req: true})
		mem := f.FlagInt(term.Flag{Name: "m", Desc: "Memory size in bits.", Req: true})
		if f.Err() != nil {
			return f.Err()
		}

		var dev *eeprom.Device
		bus := ctx.Value(ctxKeyI2C).(*i2c.I2C)
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

		_, err := mem.Seek(int64(start), 0)
		if err != nil {
			return fmt.Errorf("seek: %w", err)
		}

		r := io.Reader(mem)
		if count > 0 {
			r = io.LimitReader(r, int64(count))
		}

		wc := hex.Dumper(term.Printer(ctx))
		_, err = io.Copy(wc, r)
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}

		return wc.Close()
	}},
}
