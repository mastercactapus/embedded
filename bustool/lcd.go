package bustool

import (
	"context"
	"io"

	"github.com/mastercactapus/embedded/bus/i2c"
	"github.com/mastercactapus/embedded/driver/lcd"
	"github.com/mastercactapus/embedded/term"
)

func AddLCD(sh *term.Shell) *term.Shell {
	lcdSh := sh.NewSubShell(term.Command{Name: "lcd", Desc: "Interact with an LCD display over I2C.", Init: func(ctx context.Context, exec term.CmdFunc) error {
		f := term.Flags(ctx)
		addr := f.Uint16(term.Flag{Name: "addr", Short: 'd', Def: "0x27", Env: "DEV", Desc: "Device addresss.", Req: true})
		if err := f.Parse(); err != nil {
			return err
		}

		bus := ctx.Value(ctxKeyI2C).(i2c.Bus)
		dev := lcd.NewHD44780I2C(bus, *addr)

		err := dev.Init()
		if err != nil {
			return err
		}

		return exec(context.WithValue(ctx, ctxKeyLCD, dev))
	}})

	for _, c := range lcdCommands {
		lcdSh.AddCommand(c)
	}
	return lcdSh
}

var lcdCommands = []term.Command{
	{Name: "on", Desc: "Turn on the backlight.", Exec: func(ctx context.Context) error {
		if err := term.Flags(ctx).Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyLCD).(*lcd.HD44780)
		return dev.SetBacklight(true)
	}},
	{Name: "off", Desc: "Turn off the backlight.", Exec: func(ctx context.Context) error {
		if err := term.Flags(ctx).Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyLCD).(*lcd.HD44780)
		return dev.SetBacklight(false)
	}},
	{Name: "cls", Desc: "Clear the screen.", Exec: func(ctx context.Context) error {
		if err := term.Flags(ctx).Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyLCD).(*lcd.HD44780)
		return dev.Clear()
	}},
	{Name: "w", Desc: "Write to the screen.", Exec: func(ctx context.Context) error {
		f := term.Flags(ctx)
		x := f.Byte(term.Flag{Short: 'x', Def: "0", Desc: "Cursor start X.", Req: true})
		y := f.Byte(term.Flag{Short: 'y', Def: "0", Desc: "Cursor start Y.", Req: true})
		if err := f.Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyLCD).(*lcd.HD44780)

		err := dev.SetCursor(*y, *x)
		if err != nil {
			return err
		}

		_, err = io.WriteString(dev, f.Arg(0))
		return err
	}},
}
