package bustool

import (
	"context"
	"machine"

	"github.com/mastercactapus/embedded/driver/lcd"
	"github.com/mastercactapus/embedded/term"
	"tinygo.org/x/drivers/hd44780i2c"
)

func AddLCD(sh *term.Shell) *term.Shell {
	lcdSh := sh.NewSubShell(term.Command{Name: "lcd", Desc: "Interact with an LCD display over I2C.", Init: func(ctx context.Context, exec term.CmdFunc) error {
		f := term.Flags(ctx)
		addr := f.Byte(term.Flag{Name: "addr", Short: 'd', Def: "0x27", Env: "DEV", Desc: "Device addresss.", Req: true})
		if err := f.Parse(); err != nil {
			return err
		}

		bus := ctx.Value(ctxKeyI2C).(lcd.I2C)
		dev := lcd.NewI2C(bus, *addr)

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

		dev := ctx.Value(ctxKeyLCD).(*lcd.Device)
		return dev.SetBacklight(true)
	}},
	{Name: "off", Desc: "Turn off the backlight.", Exec: func(ctx context.Context) error {
		if err := term.Flags(ctx).Parse(); err != nil {
			return err
		}

		dev := ctx.Value(ctxKeyLCD).(*lcd.Device)
		return dev.SetBacklight(false)
	}},
	{Name: "other", Desc: "Test other init.", Exec: func(ctx context.Context) error {
		if err := term.Flags(ctx).Parse(); err != nil {
			return err
		}

		err := machine.I2C0.Configure(machine.I2CConfig{})
		if err != nil {
			return err
		}

		lcd := hd44780i2c.New(machine.I2C0, 0x27)
		err = lcd.Configure(hd44780i2c.Config{
			Height:      1,
			Width:       8,
			CursorOn:    true,
			CursorBlink: true,
		})
		if err != nil {
			return err
		}
		lcd.Print([]byte("H"))

		return nil
	}},
}
