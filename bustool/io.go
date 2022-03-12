package bustool

import (
	"fmt"
	"strconv"

	"github.com/mastercactapus/embedded/driver/ioexp"
	"github.com/mastercactapus/embedded/serial/i2c"
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
	return ioexp.PinMask(pins...), nil
}

func AddIO(sh *term.Shell) *term.Shell {
	ioSh := sh.NewSubShell("io", "Interact with IO expansion chips over I2C.", func(r term.RunArgs) error {
		addr := r.Uint16(term.Flag{Name: "addr", Short: 'd', Def: "0x20", Env: "DEV", Desc: "Device addresss.", Req: true})
		devType := r.Enum(term.Flag{Name: "type", Short: 't', Def: "mcp16", Desc: "IO device type.", Req: true}, "pcf", "mcp8", "mcp16")
		if err := r.Parse(); err != nil {
			return err
		}

		var dev ioexp.PinReadWriter
		bus := r.Get("i2c").(i2c.Bus)

		switch *devType {
		case "mcp16":
			dev = ioexp.NewMCP23017(bus, *addr)
		case "mcp8":
			dev = ioexp.NewMCP23008(bus, *addr)
		case "pcf":
			dev = ioexp.NewPCF8574(bus, *addr)
		default:
			return r.UsageError("unsupported device type '%s'", *devType)
		}
		r.Set("io", dev)

		return nil
	})

	ioSh.AddCommands(ioCommands...)
	return ioSh
}

var ioCommands = []term.Command{
	{Name: "r", Desc: "Read pin state.", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("io").(ioexp.PinReadWriter)
		pins, err := dev.ReadPins()
		if err != nil {
			return err
		}

		for i := 0; i < dev.PinCount(); i++ {
			r.Printf("% 3d ", i)
		}
		r.Println()
		for i := 0; i < dev.PinCount(); i++ {
			val := 0
			if pins.Value(i) {
				val = 1
			}
			r.Printf("% 3d ", val)
		}
		r.Println()

		return nil
	}},

	{Name: "input", Desc: "Set selected pins to input, rest as output.", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("io").(ioexp.PinReadWriter)
		is, ok := dev.(ioexp.InputSetter)
		if !ok {
			return fmt.Errorf("device does not support setting input pins")
		}

		mask, err := pinMask(r.Args())
		if err != nil {
			return r.UsageError("parse args: %w", err)
		}

		return is.SetInputPinsMask(ioexp.AllPins(true), mask)
	}},

	{Name: "output", Desc: "Set selected pins to output, rest as input.", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("io").(ioexp.PinReadWriter)
		is, ok := dev.(ioexp.InputSetter)
		if !ok {
			return fmt.Errorf("device does not support setting input pins")
		}

		mask, err := pinMask(r.Args())
		if err != nil {
			return r.UsageError("parse args: %w", err)
		}

		return is.SetInputPinsMask(ioexp.AllPins(false), mask)
	}},

	{Name: "high", Desc: "Turn HIGH selected pin(s).", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("io").(ioexp.PinReadWriter)

		mask, err := pinMask(r.Args())
		if err != nil {
			return r.UsageError("parse args: %w", err)
		}

		return dev.WritePinsMask(ioexp.AllPins(true), mask)
	}},

	{Name: "low", Desc: "Turn LOW selected pin(s).", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("io").(ioexp.PinReadWriter)

		mask, err := pinMask(r.Args())
		if err != nil {
			return r.UsageError("parse args: %w", err)
		}

		return dev.WritePinsMask(ioexp.AllPins(false), mask)
	}},

	{Name: "set", Desc: "Set specified pins.", Exec: func(r term.RunArgs) error {
		low := r.Bool(term.Flag{Name: "low", Short: 'l', Desc: "Set specified pins LOW instead of HIGH."})
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("io").(ioexp.PinReadWriter)

		pins, err := pinMask(r.Args())
		if err != nil {
			return r.UsageError("parse args: %w", err)
		}

		if *low {
			pins = ioexp.Invert(pins)
		}

		return dev.WritePins(pins)
	}},
}
