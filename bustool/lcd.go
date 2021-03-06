package bustool

import (
	"io"
	"math/rand"

	"github.com/mastercactapus/embedded/driver/lcd"
	"github.com/mastercactapus/embedded/serial/i2c"
	"github.com/mastercactapus/embedded/term"
)

func AddLCD(sh *term.Shell) *term.Shell {
	lcdSh := sh.NewSubShell("lcd", "Interact with an LCD display over I2C.", func(r term.RunArgs) error {
		addr := r.Uint16(term.Flag{Name: "addr", Short: 'd', Def: "0x27", Env: "DEV", Desc: "Device addresss.", Req: true})
		lines := r.Int(term.Flag{Name: "lines", Short: 'l', Def: "2", Env: "LINES", Desc: "Number of lines.", Req: true})
		cols := r.Int(term.Flag{Name: "cols", Short: 'c', Def: "16", Env: "COLS", Desc: "Number of columns.", Req: true})
		if err := r.Parse(); err != nil {
			return err
		}

		bus := r.Get("i2c").(i2c.Bus)
		dev, err := lcd.NewHD44780I2C(bus, *addr, lcd.Config{
			Lines: *lines,
			Cols:  *cols,
		})
		if err != nil {
			return err
		}
		r.Set("lcd", dev)

		return nil
	})

	lcdSh.AddCommands(LCDCommands...)
	return lcdSh
}

// LCDCommands are commands for interacting with an HD44780 display.
//
// The device must be available at the 'lcd' key.
var LCDCommands = []term.Command{
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
	{Name: "type", Desc: "Type a string on the screen.", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("lcd").(*lcd.HD44780)
		for {
			c := <-r.Input()
			if c == 0 {
				continue
			}
			if c == term.Interrupt {
				return nil
			}
			dev.WriteByte(byte(c))
		}
	}},
	{Name: "stress", Desc: "Stream data to the display.", Exec: func(r term.RunArgs) error {
		randData := r.Bool(term.Flag{Name: "rand", Short: 'r', Desc: "Randomize the data."})
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("lcd").(*lcd.HD44780)

		d := byte('a')

		for r.WaitForInterrupt() {
			if *randData {
				dev.WriteByte(byte(rand.Intn('~'-' ')) + ' ')
			} else {
				dev.WriteByte(d)
				d++
				if d > 'z' {
					d = 'a'
				}
			}
		}
		return nil
	}},
	{Name: "w", Desc: "Write to the screen.", Exec: func(r term.RunArgs) error {
		x := r.Int(term.Flag{Short: 'x', Def: "0", Desc: "Cursor start X.", Req: true})
		y := r.Int(term.Flag{Short: 'y', Def: "0", Desc: "Cursor start Y.", Req: true})
		if err := r.Parse(); err != nil {
			return err
		}

		dev := r.Get("lcd").(*lcd.HD44780)

		err := dev.SetCursorXY(*x, *y)
		if err != nil {
			return err
		}

		_, err = io.WriteString(dev, r.Arg(0))
		return err
	}},
}
