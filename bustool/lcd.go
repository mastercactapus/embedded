package bustool

import (
	"io"

	"github.com/mastercactapus/embedded/bus/i2c"
	"github.com/mastercactapus/embedded/driver/lcd"
	"github.com/mastercactapus/embedded/term"
)

func AddLCD(sh *term.Shell) *term.Shell {
	lcdSh := sh.NewSubShell(term.Command{Name: "lcd", Desc: "Interact with an LCD display over I2C.", Exec: func(r term.RunArgs) error {
		addr := r.Uint16(term.Flag{Name: "addr", Short: 'd', Def: "0x27", Env: "DEV", Desc: "Device addresss.", Req: true})
		if err := r.Parse(); err != nil {
			return err
		}

		bus := r.Get("i2c").(i2c.Bus)
		dev := lcd.NewHD44780I2C(bus, *addr)

		err := dev.Init()
		if err != nil {
			return err
		}
		r.Set("lcd", dev)

		return nil
	}})

	for _, c := range lcdCommands {
		lcdSh.AddCommand(c)
	}
	return lcdSh
}

var lcdCommands = []term.Command{
	{Name: "on", Desc: "Turn on the backlight.", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("lcd").(*lcd.HD44780)
		return dev.SetBacklight(true)
	}},
	{Name: "off", Desc: "Turn off the backlight.", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("lcd").(*lcd.HD44780)
		return dev.SetBacklight(false)
	}},
	{Name: "cls", Desc: "Clear the screen.", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("lcd").(*lcd.HD44780)
		return dev.Clear()
	}},
	{Name: "w", Desc: "Write to the screen.", Exec: func(r term.RunArgs) error {
		x := r.Byte(term.Flag{Short: 'x', Def: "0", Desc: "Cursor start X.", Req: true})
		y := r.Byte(term.Flag{Short: 'y', Def: "0", Desc: "Cursor start Y.", Req: true})
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("lcd").(*lcd.HD44780)

		err := dev.SetCursor(*x, *y)
		if err != nil {
			return err
		}

		_, err = io.WriteString(dev, r.Arg(0))
		return err
	}},
}
