package bustool

import (
	"io"

	"github.com/mastercactapus/embedded/i2c"
	"github.com/mastercactapus/embedded/term"
)

func NewShell(r io.Reader, w io.Writer) *term.Shell {
	sh := term.NewShell("bustool", "Interact with various embedded devices.", r, w)
	sh.AddCommand(term.Command{Name: "version", Desc: "Output version information.", Exec: func(c *term.CmdCtx) error {
		c.Printer().Println("v0")

		return nil
	}})

	i2cSh := sh.NewSubShell(term.Command{Name: "i2c", Desc: "Interact with I2C devices.", Exec: func(c *term.CmdCtx) error {
		bus := i2c.New()
		bus.Configure(i2c.Config{
			// SCL: i2cPin(machine.I2C0_SCL_PIN),
			// SDA: i2cPin(machine.I2C0_SDA_PIN),
		})
		c.Set("i2c", bus)
		return nil
	}})

	i2cSh.AddCommand(term.Command{Name: "ping", Desc: "Ping a device.", Exec: func(c *term.CmdCtx) error {
		addr := c.FlagByte(term.Flag{Name: "d", Env: "DEVICE", Desc: "Device addresss.", Req: true})
		write := c.FlagBool(term.Flag{Name: "w", Desc: "Ping the write address instead."})
		c.Parse()

		bus := c.Get("i2c").(*i2c.I2C)

		if write {
			return bus.PingW(addr)
		}

		return bus.Ping(addr)
	}})

	return sh
}
