package term

import (
	"io"
	"strings"

	"github.com/mastercactapus/embedded/term/ansi"
)

func (fp *FlagParser) PrintUsage(w io.Writer) {
	// cmd.addFlag(Flag{Name: "h", Desc: "Show this help message"}, "")
	// p := cmd.Printer()
	// if err != nil {
	// 	p.Fg(ansi.Red)
	// 	p.Println(err.Error())
	// 	p.Reset()
	// }

	// p.Println(fp.c.Desc)
	// p.Println()

	p := ansi.NewPrinter(w)
	p.Printf("Usage: %s", fp.cmd.Name)
	if len(fp.flagInfo) > 0 {
		p.Printf(" [flags]")
	}
	for _, a := range fp.argInfo {
		switch {
		case a.isSlice && a.Req:
			p.Printf(" <%s ...>", a.Arg.Name)
		case a.isSlice && !a.Req:
			p.Printf(" [%s ...]", a.Arg.Name)
		case !a.isSlice && a.Req:
			p.Printf(" <%s>", a.Arg.Name)
		case !a.isSlice && !a.Req:
			p.Printf(" [%s]", a.Arg.Name)
		}
	}
	p.Println()
	p.Println()

	if len(fp.examples) > 0 {
		p.Println("Examples:")
		for _, ex := range fp.examples {
			p.Printf("    %s\n", ex.cmdline)
			p.Printf("        %s\n", ex.details)
		}
		p.Println()
	}

	var tb ansi.Table
	tb.Pad = 4

	if len(fp.argInfo) > 0 {
		tb.AddLine("Arguments:")
		for _, a := range fp.argInfo {
			desc := a.Desc
			if a.Req {
				desc = "REQUIRED: " + desc
			}
			tb.AddRow("    "+a.Name, a.typeName, desc, "")
		}
		tb.AddLine("")
	}

	if len(fp.flagInfo) > 0 {
		tb.AddLine("Flags:")
		for _, f := range fp.flagInfo {
			env := f.Env
			if env != "" {
				env = "[$" + env + "]"
			}
			desc := f.Desc
			if f.Def != "" {
				desc += " (default: " + f.Def + ")"
			}
			if f.Req {
				desc = "REQUIRED: " + desc
			}
			tb.AddRow("    -"+f.Name, f.typeName, desc, env)
		}
		tb.AddLine("")
	}

	p.Println(strings.TrimSpace(tb.String()))
}
