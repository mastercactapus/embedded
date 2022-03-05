package bustool

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mastercactapus/embedded/bus/i2c"
	"github.com/mastercactapus/embedded/driver/ioexp"
	"github.com/mastercactapus/embedded/term"
)

func pinMask(args []string) (ioexp.Valuer, error) {
	var pins []int
	for _, arg := range args {
		if arg == "all" {
			return ioexp.AllPins(true), nil
		}
		i, err := strconv.Atoi(arg)
		if err != nil {
			return nil, err
		}
		pins = append(pins, i)
	}
	return ioexp.PinMask(pins), nil
}

func AddIO(sh *term.Shell) *term.Shell {
	ioSh := sh.NewSubShell(term.Command{Name: "io", Desc: "Interact with IO expansion chips over I2C.", Init: func(ctx context.Context, exec term.CmdFunc) error {
		f := term.Flags(ctx)
		addr := f.Uint16(term.Flag{Name: "addr", Short: 'd', Def: "0x20", Env: "DEV", Desc: "Device addresss.", Req: true})
		devType := f.Enum(term.Flag{Name: "type", Short: 't', Def: "mcp", Desc: "IO device type.", Req: true}, "pcf", "mcp")
		if err := f.Parse(); err != nil {
			return err
		}

		var dev ioexp.PinReadWriter
		bus := ctx.Value(ctxKeyI2C).(i2c.Bus)

		switch *devType {
		case "mcp":
			dev = ioexp.NewMCP23017(bus, *addr)
		case "pcf":
			dev = ioexp.NewPCF8574(bus, *addr)
		default:
			return f.UsageError("unsupported device type '%s'", *devType)
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

		p := term.Printer(ctx)
		for i := 0; i < dev.PinCount(); i++ {
			p.Printf("% 3d ", i)
		}
		p.Println()
		for i := 0; i < dev.PinCount(); i++ {
			val := 0
			if pins.Value(i) {
				val = 1
			}
			p.Printf("% 3d ", val)
		}
		p.Println()

		return nil
	}},

	{Name: "input", Desc: "Set selected pins to input, rest as output.", Exec: func(ctx context.Context) error {
		f := term.Flags(ctx)
		if err := f.Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyIO).(ioexp.PinReadWriter)
		is, ok := dev.(ioexp.InputSetter)
		if !ok {
			return fmt.Errorf("device does not support setting input pins")
		}

		mask, err := pinMask(f.Args())
		if err != nil {
			return f.UsageError("parse args: %w", err)
		}

		return is.SetInputPinsMask(ioexp.AllPins(true), mask)
	}},

	{Name: "output", Desc: "Set selected pins to output, rest as input.", Exec: func(ctx context.Context) error {
		f := term.Flags(ctx)
		if err := f.Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyIO).(ioexp.PinReadWriter)
		is, ok := dev.(ioexp.InputSetter)
		if !ok {
			return fmt.Errorf("device does not support setting input pins")
		}

		mask, err := pinMask(f.Args())
		if err != nil {
			return f.UsageError("parse args: %w", err)
		}

		return is.SetInputPinsMask(ioexp.AllPins(false), mask)
	}},

	{Name: "on", Desc: "Turn on selected pin(s).", Exec: func(ctx context.Context) error {
		f := term.Flags(ctx)
		if err := f.Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyIO).(ioexp.PinReadWriter)

		mask, err := pinMask(f.Args())
		if err != nil {
			return f.UsageError("parse args: %w", err)
		}

		return dev.WritePinsMask(ioexp.AllPins(true), mask)
	}},

	{Name: "off", Desc: "Turn off selected pin(s).", Exec: func(ctx context.Context) error {
		f := term.Flags(ctx)
		if err := f.Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyIO).(ioexp.PinReadWriter)

		mask, err := pinMask(f.Args())
		if err != nil {
			return f.UsageError("parse args: %w", err)
		}

		return dev.WritePinsMask(ioexp.AllPins(false), mask)
	}},

	{Name: "set", Desc: "Turn on ONLY selected pin(s).", Exec: func(ctx context.Context) error {
		f := term.Flags(ctx)
		if err := f.Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyIO).(ioexp.PinReadWriter)

		mask, err := pinMask(f.Args())
		if err != nil {
			return f.UsageError("parse args: %w", err)
		}

		return dev.WritePins(mask)
	}},
}
