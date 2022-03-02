package term

import (
	"context"
	"errors"
	"strings"

	"github.com/mastercactapus/embedded/term/ansi"
)

var builtins = []Command{
	{"help", "Print this help message.", func(ctx context.Context) error {
		if err := Parse(ctx).Err(); err != nil {
			return err
		}
		p := Printer(ctx)

		p.Reset()

		cmd := ctx.Value(ctxKeyCmd).(*cmdContext)
		if cmd.desc != "" {
			p.Println(cmd.desc)
			p.Println()
		}

		sh := cmd.sh

		var tb ansi.Table
		tb.Min = 3

		tb.AddLine("Built-in:")
		for _, name := range sh.bNames {
			tb.AddRow("    "+name, sh.cmds[name].Desc)
		}
		tb.AddLine("")

		if len(sh.shNames) > 0 {
			tb.AddLine("Sub-shells:")
			for _, name := range sh.shNames {
				tb.AddRow("    "+name, sh.cmds[name].Desc)
			}
			tb.AddLine("")
		}

		if len(sh.cmdNames) > 0 {
			tb.AddLine("Commands:")
			for _, name := range sh.cmdNames {
				tb.AddRow("    "+name, sh.cmds[name].Desc)
			}
			tb.AddLine("")
		}

		p.Print(tb.String())

		return nil
	}, nil},

	{"clear", "Clear the screen.", func(ctx context.Context) error {
		if err := Parse(ctx).Err(); err != nil {
			return err
		}

		Printer(ctx).Esc('J', 2)
		return nil
	}, nil},

	{"export", "Export envronment variables.", func(ctx context.Context) error {
		f := Parse(ctx)
		f.Example("export FOO=1", "Set the FOO variable to 1.")
		f.Example("export FOO=", "Clear the FOO variable.")
		flags := f.ArgStringN(Arg{Name: "ENV", Desc: "KEY=value pairs of env variables to set.", Req: true})
		if err := f.Err(); err != nil {
			return err
		}

		sh := ctx.Value(ctxKeyCmd).(*cmdContext).sh

		toSet := flags[:0]
		for _, flag := range flags {
			name, value := ansi.Cut(flag, '=')
			if value != "" {
				toSet = append(toSet, flag)
			}
			prefix := name + "="
			for i, pair := range sh.env {
				if !strings.HasPrefix(pair, prefix) {
					continue
				}
				// remove the old one
				sh.env = append(sh.env[:i], sh.env[i+1:]...)
				break
			}
		}

		sh.env = append(sh.env, toSet...)

		return nil
	}, nil},

	{"env", "Print shell environment values.", func(ctx context.Context) error {
		if err := Parse(ctx).Err(); err != nil {
			return err
		}

		cmd := ctx.Value(ctxKeyCmd).(*cmdContext)
		p := Printer(ctx)
		for _, pair := range cmd.sh.env {
			p.Println(pair)
		}

		return nil
	}, nil},

	{"reset", "Reset all environment variables.", func(ctx context.Context) error {
		if err := Parse(ctx).Err(); err != nil {
			return err
		}

		cmd := ctx.Value(ctxKeyCmd).(*cmdContext)
		cmd.sh.env = cmd.sh.env[:0]
		return nil
	}, nil},

	{"exit", "Exits the current shell.", func(ctx context.Context) error {
		f := Parse(ctx)
		errMsg := f.ArgString(Arg{Name: "message", Desc: "If not empty, exits the shell with error message."})
		if err := f.Err(); err != nil {
			return err
		}

		if errMsg != "" {
			return exitErr{errors.New(errMsg)}
		}

		return exitErr{nil}
	}, nil},
}
