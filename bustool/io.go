package bustool

import (
	"strconv"

	"github.com/mastercactapus/embedded/driver"
	"github.com/mastercactapus/embedded/driver/ioexp"
	"github.com/mastercactapus/embedded/serial/i2c"
	"github.com/mastercactapus/embedded/term"
)

func AddIO(sh *term.Shell) *term.Shell {
	ioSh := sh.NewSubShell("io", "Interact with IO expansion chips over I2C.", func(r term.RunArgs) error {
		addr := r.Uint16(term.Flag{Name: "addr", Short: 'd', Def: "0x20", Env: "DEV", Desc: "Device addresss.", Req: true})
		devType := r.Enum(term.Flag{Name: "type", Short: 't', Def: "mcp16", Desc: "IO device type.", Req: true}, "pcf", "mcp8", "mcp16")
		if err := r.Parse(); err != nil {
			return err
		}

		var dev driver.Pinner
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

		dev := r.Get("io").(driver.Pinner)

		for i := 0; i < dev.PinCount(); i++ {
			r.Printf("% 3d ", i)
		}
		r.Println()
		for i := 0; i < dev.PinCount(); i++ {
			val := 0
			v, err := dev.Pin(i).Get()
			if err != nil {
				return err
			}
			if v {
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

		dev := r.Get("io").(driver.Pinner)
		for _, arg := range r.Args() {
			i, err := strconv.Atoi(arg)
			if err != nil {
				return err
			}
			if err := dev.Pin(i).Input(); err != nil {
				return err
			}
		}

		return nil
	}},

	{Name: "output", Desc: "Set selected pins to output, rest as input.", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("io").(driver.Pinner)
		for _, arg := range r.Args() {
			i, err := strconv.Atoi(arg)
			if err != nil {
				return err
			}
			if err := dev.Pin(i).Output(); err != nil {
				return err
			}
		}

		return nil
	}},

	{Name: "high", Desc: "Turn HIGH selected pin(s).", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("io").(driver.Pinner)
		for _, arg := range r.Args() {
			i, err := strconv.Atoi(arg)
			if err != nil {
				return err
			}
			if err := dev.Pin(i).High(); err != nil {
				return err
			}
		}

		return nil
	}},

	{Name: "low", Desc: "Turn LOW selected pin(s).", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("io").(driver.Pinner)
		for _, arg := range r.Args() {
			i, err := strconv.Atoi(arg)
			if err != nil {
				return err
			}
			if err := dev.Pin(i).Low(); err != nil {
				return err
			}
		}

		return nil
	}},
}
