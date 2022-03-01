package term

import (
	"errors"
	"strings"

	"github.com/mastercactapus/embedded/term/ansi"
)

var builtins = []Command{
	{"help", "Print this help message.", func(c *CmdCtx) error {
		c.Parse()
		p := c.Printer()

		p.Reset()
		if c.c.sh.desc != "" {
			p.Println(c.c.sh.desc)
			p.Println()
		}

		var tb ansi.Table
		tb.Min = 3

		tb.AddLine("Built-in:")
		for _, name := range c.c.sh.bNames {
			tb.AddRow("    "+name, c.c.sh.cmds[name].Desc)
		}
		tb.AddLine("")

		if len(c.c.sh.shNames) > 0 {
			tb.AddLine("Sub-shells:")
			for _, name := range c.c.sh.shNames {
				tb.AddRow("    "+name, c.c.sh.cmds[name].Desc)
			}
			tb.AddLine("")
		}

		if len(c.c.sh.cmdNames) > 0 {
			tb.AddLine("Commands:")
			for _, name := range c.c.sh.cmdNames {
				tb.AddRow("    "+name, c.c.sh.cmds[name].Desc)
			}
			tb.AddLine("")
		}

		p.Print(tb.String())

		return nil
	}},

	{"clear", "Clear the screen.", func(c *CmdCtx) error {
		c.Parse()
		c.Printer().Esc('c')
		return nil
	}},

	{"export", "Export envronment variables.", func(c *CmdCtx) error {
		c.Example("export FOO=1", "Set the FOO variable to 1.")
		c.Example("export FOO=", "Clear the FOO variable.")
		flags := c.ArgStringN(Arg{Name: "ENV", Desc: "KEY=value pairs of env variables to set.", Req: true})
		c.Parse()

		toSet := flags[:0]
		for _, flag := range flags {
			name, value := ansi.Cut(flag, '=')
			if value != "" {
				toSet = append(toSet, flag)
			}
			prefix := name + "="
			for i, pair := range c.c.sh.env {
				if !strings.HasPrefix(pair, prefix) {
					continue
				}
				// remove the old one
				c.c.sh.env = append(c.c.sh.env[:i], c.c.sh.env[i+1:]...)
			}
		}

		c.c.sh.env = append(c.c.sh.env, toSet...)

		return nil
	}},

	{"env", "Print environment values.", func(c *CmdCtx) error {
		c.Parse()

		for _, pair := range c.c.sh.env {
			c.Printer().Println(pair)
		}

		return nil
	}},

	{"reset", "Reset all environment variables.", func(c *CmdCtx) error {
		c.Parse()

		c.c.sh.env = c.c.sh.env[:0]
		return nil
	}},

	{"exit", "Exits the current shell.", func(c *CmdCtx) error {
		errMsg := c.ArgString(Arg{Name: "message", Desc: "If not empty, exits the shell with error message."})
		c.Parse()

		if errMsg != "" {
			return exitErr{errors.New(errMsg)}
		}

		return exitErr{nil}
	}},
}
