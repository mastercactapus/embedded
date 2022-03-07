package bustool

import (
	"io"

	"github.com/mastercactapus/embedded/term"
)

func NewShell(r io.Reader, w io.Writer) *term.Shell {
	sh := term.NewRootShell("bustool", "Interact with various embedded devices.", r, w)
	sh.AddCommand(term.Command{Name: "version", Desc: "Output version information.", Exec: func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}
		r.Println("v1")

		return nil
	}})

	return sh
}
