package term

import (
	"strings"

	"github.com/mastercactapus/embedded/term/ansi"
)

func (cmd *CmdCtx) usage(err error) {
	cmd.addFlag(Flag{Name: "h", Desc: "Show this help message"}, "")
	p := cmd.Printer()
	if err != nil {
		p.Fg(ansi.Red)
		p.Println(err.Error())
		p.Reset()
	}

	p.Println(cmd.c.Desc)
	p.Println()

	p.Printf("Usage: %s", cmd.c.Name)
	if len(cmd.flagInfo) > 0 {
		p.Printf(" [flags]")
	}
	for _, a := range cmd.argInfo {
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

	if len(cmd.examples) > 0 {
		p.Println("Examples:")
		for _, ex := range cmd.examples {
			p.Printf("    %s\n", ex.cmdline)
			p.Printf("        %s\n", ex.details)
		}
		p.Println()
	}

	var tb ansi.Table
	tb.Pad = 4

	if len(cmd.argInfo) > 0 {
		tb.AddLine("Arguments:")
		for _, a := range cmd.argInfo {
			desc := a.Desc
			if a.Req {
				desc = "REQUIRED: " + desc
			}
			tb.AddRow("    "+a.Name, a.typeName, desc, "")
		}
		tb.AddLine("")
	}

	if len(cmd.flagInfo) > 0 {
		tb.AddLine("Flags:")
		for _, f := range cmd.flagInfo {
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
