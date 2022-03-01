package bustool

import (
	"io"

	"github.com/mastercactapus/embedded/term"
)

func NewShell(r io.Reader, w io.Writer) *term.Shell {
	sh := term.NewShell("bustool", "Interact with various embedded devices.", r, w)
	sh.AddCommand(term.Command{Name: "version", Desc: "Output version information.", Exec: func(c *term.CmdCtx) error {
		if err := c.Parse(); err != nil {
			return err
		}
		c.Printer().Println("v0")

		return nil
	}})

	return sh
}
