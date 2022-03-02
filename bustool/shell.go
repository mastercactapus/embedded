package bustool

import (
	"context"
	"io"

	"github.com/mastercactapus/embedded/term"
)

func NewShell(r io.Reader, w io.Writer) *term.Shell {
	sh := term.NewShell("bustool", "Interact with various embedded devices.", r, w)
	sh.AddCommand(term.Command{Name: "version", Desc: "Output version information.", Exec: func(ctx context.Context) error {
		if err := term.Parse(ctx).Err(); err != nil {
			return err
		}
		term.Printer(ctx).Println("v0")

		return nil
	}})

	return sh
}
