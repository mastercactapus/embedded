package term

import (
	"context"
	"errors"
	"strings"

	"github.com/mastercactapus/embedded/term/ansi"
)

var builtins = []Command{
	{"help", "Print this help message.", func(ctx context.Context) error {
		if err := Flags(ctx).Parse(); err != nil {
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
		if err := Flags(ctx).Parse(); err != nil {
			return err
		}

		Printer(ctx).Esc('J', 2)
		return nil
	}, nil},

	{"export", "Export envronment variables.", func(ctx context.Context) error {
		f := Flags(ctx)
		f.SetHelpParameters("[name=[value]] ...]")
		f.Example("export FOO=1", "Set the FOO variable to 1.")
		f.Example("export FOO=", "Clear the FOO variable.")
		unset := f.Bool(Flag{Short: 'n', Desc: "Remove variables from the environment."})
		if err := f.Parse(); err != nil {
			return err
		}

		// TODO: usage and parsing export Args
		cmd := ctx.Value(ctxKeyCmd).(*cmdContext)
		if *unset {
			for _, arg := range f.set.Args() {
				name, _ := ansi.Cut(arg, '=')
				cmd.env.Unset(name)
			}
			return nil
		}

		for _, arg := range f.set.Args() {
			if !strings.Contains(arg, "=") {
				val, ok := cmd.env.Get(arg)
				if ok {
					cmd.env.Set(arg, val)
				}
				continue
			}
			cmd.env.Set(ansi.Cut(arg, '='))
		}

		return nil
	}, nil},

	{"env", "Print shell environment values.", func(ctx context.Context) error {
		if err := Flags(ctx).Parse(); err != nil {
			return err
		}

		cmd := ctx.Value(ctxKeyCmd).(*cmdContext)
		p := Printer(ctx)
		for _, key := range cmd.env.List() {
			val, _ := cmd.env.Get(key)
			p.Println(key + "=" + val)
		}

		return nil
	}, nil},

	{"exit", "Exits the current shell.", func(ctx context.Context) error {
		f := Flags(ctx)
		errMsg := f.String(Flag{Short: 'm', Desc: "Exit shell with an error message."})
		if err := f.Parse(); err != nil {
			return err
		}

		if errMsg != nil {
			return exitErr{errors.New(*errMsg)}
		}

		return exitErr{nil}
	}, nil},
}
