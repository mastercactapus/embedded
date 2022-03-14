package term

import (
	"io"

	"github.com/mastercactapus/embedded/term/ansi"
)

func (e usageErr) PrintUsage(w io.Writer) {
	p := ansi.NewPrinter(w)

	p.Printf("%s [flags ...] %s\n", e.fs.cmd.Args[0], e.fs.helpParams)
	p.Println()

	for _, id := range e.fs.flagList {
		f := e.fs.flags[id]
		switch {
		case f.Short != 0 && f.Name != "":
			p.Printf("\t-%c, --%s", f.Short, f.Name)
		case f.Short != 0:
			p.Printf("\t-%c", f.Short)
		case f.Name != "":
			p.Printf("\t--%s", f.Name)
		}
		if f.Type != "" {
			p.Printf(" <%s>", f.Type)
		}
		p.Println()
		p.Print("\t\t")
		if f.Req {
			p.Print("Required. ")
		}
		p.Println(f.Desc)
		if f.Def != "" {
			p.Printf("(default: %s)\n", f.Def)
		}
	}
}
