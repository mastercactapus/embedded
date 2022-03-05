package bustool

import (
	"context"
	"strconv"

	"github.com/mastercactapus/embedded/driver/ioexp"
	"github.com/mastercactapus/embedded/i2c"
	"github.com/mastercactapus/embedded/term"
	"github.com/mastercactapus/embedded/term/ansi"
)

func AddIO(sh *term.Shell) *term.Shell {
	ioSh := sh.NewSubShell(term.Command{Name: "io", Desc: "Interact with IO expansion chips over I2C.", Init: func(ctx context.Context, exec term.CmdFunc) error {
		f := term.Flags(ctx)
		addr := f.Uint16(term.Flag{Name: "addr", Short: 'd', Def: "0x20", Env: "DEV", Desc: "Device addresss.", Req: true})
		pinN := f.Int(term.Flag{Name: "pins", Short: 'p', Def: "8", Desc: "Pin count.", Req: true})
		if err := f.Parse(); err != nil {
			return err
		}

		var dev ioexp.PinReadWriter
		bus := ctx.Value(ctxKeyI2C).(i2c.Bus)

		switch *pinN {
		case 8:
			dev = ioexp.NewPCF8574(bus, *addr)
		default:
			return f.UsageError("unsupported pin count %d", pinN)
		}

		return exec(context.WithValue(ctx, ctxKeyIO, dev))
	}})

	for _, c := range ioCommands {
		ioSh.AddCommand(c)
	}
	return ioSh
}

var ioCommands = []term.Command{
	{Name: "r", Desc: "Read pin state.", Exec: func(ctx context.Context) error {
		if err := term.Flags(ctx).Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyIO).(ioexp.PinReadWriter)
		pins, err := dev.ReadPins()
		if err != nil {
			return err
		}

		var t ansi.Table
		t.AddRow("0", "1", "2", "3", "4", "5", "6", "7")
		var state []string
		for i := 0; i < 8; i++ {
			if pins.Value(i) {
				state = append(state, "H")
			} else {
				state = append(state, "L")
			}
		}
		t.AddRow(state...)

		term.Printer(ctx).Println(t.String())

		return nil
	}},

	{Name: "on", Desc: "Turn on selected pin(s).", Exec: func(ctx context.Context) error {
		if err := term.Flags(ctx).Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyIO).(ioexp.PinReadWriter)
		pins, err := dev.ReadPins()
		if err != nil {
			return err
		}
		for _, a := range term.Flags(ctx).Args() {
			if a == "all" {
				pins.SetAll(true)
				break
			}
			n, err := strconv.Atoi(a)
			if err != nil {
				return err
			}

			pins.Set(n, true)
		}
		return dev.WritePins(pins)
	}},

	{Name: "off", Desc: "Turn off selected pin(s).", Exec: func(ctx context.Context) error {
		if err := term.Flags(ctx).Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyIO).(ioexp.PinReadWriter)
		pins, err := dev.ReadPins()
		if err != nil {
			return err
		}
		for _, a := range term.Flags(ctx).Args() {
			if a == "all" {
				pins.SetAll(false)
				break
			}
			n, err := strconv.Atoi(a)
			if err != nil {
				return err
			}

			pins.Set(n, false)
		}
		return dev.WritePins(pins)
	}},

	{Name: "set", Desc: "Turn on ONLY selected pin(s).", Exec: func(ctx context.Context) error {
		if err := term.Flags(ctx).Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyIO).(ioexp.PinReadWriter)
		var pins ioexp.Pin8
		for _, a := range term.Flags(ctx).Args() {
			if a == "all" {
				pins.SetAll(true)
				break
			}
			n, err := strconv.Atoi(a)
			if err != nil {
				return err
			}

			pins.Set(n, true)
		}
		return dev.WritePins(pins)
	}},

	{Name: "toggle", Desc: "Toggle selected pin(s).", Exec: func(ctx context.Context) error {
		if err := term.Flags(ctx).Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyIO).(ioexp.PinReadWriter)
		pins, err := dev.ReadPins()
		if err != nil {
			return err
		}
		for _, a := range term.Flags(ctx).Args() {
			if a == "all" {
				pins.ToggleAll()
				continue
			}
			n, err := strconv.Atoi(a)
			if err != nil {
				return err
			}

			pins.Toggle(n)
		}
		return dev.WritePins(pins)
	}},
}
