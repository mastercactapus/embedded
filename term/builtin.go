package term

import (
	"errors"
	"sort"
	"strings"

	"github.com/mastercactapus/embedded/term/ansi"
)

var builtin []Command

func init() {
	builtin = append(builtin,
		Command{Name: "help", Desc: "Display this help message", Exec: func(r RunArgs) error {
			if err := r.Parse(); err != nil {
				return err
			}

			if r.sh.desc != "" {
				r.Println(r.sh.desc)
				r.Println()
			}

			var tb ansi.Table
			tb.Min = 3

			tb.AddLine("Built-in:")
			for _, cmd := range builtin {
				if r.sh.noExit && cmd.Name == "exit" {
					continue
				}

				tb.AddRow("    "+cmd.Name, cmd.Desc)
			}
			tb.AddLine("")

			if len(r.sh.shells) > 0 {
				sort.Slice(r.sh.shells, func(i, j int) bool {
					return r.sh.shells[i].Name < r.sh.shells[j].Name
				})
				tb.AddLine("Sub-shells:")
				for _, cmd := range r.sh.shells {
					tb.AddRow("    "+cmd.Name, cmd.Desc)
				}
				tb.AddLine("")
			}

			if len(r.sh.commands) > 0 {
				sort.Slice(r.sh.commands, func(i, j int) bool {
					return r.sh.commands[i].Name < r.sh.commands[j].Name
				})
				tb.AddLine("Commands:")
				for _, cmd := range r.sh.commands {
					tb.AddRow("    "+cmd.Name, cmd.Desc)
				}
				tb.AddLine("")
			}

			r.Print(tb.String())
			return nil
		}},
		Command{Name: "clear", Desc: "Clear the screen.", Exec: func(r RunArgs) error {
			if err := r.Parse(); err != nil {
				return err
			}

			r.Esc('J', 2)
			return nil
		}},
		Command{Name: "export", Desc: "Export envronment variables.", Exec: func(r RunArgs) error {
			r.SetHelpParameters("[name=[value]] ...]")
			r.Example("export FOO=1", "Set the FOO variable to 1.")
			r.Example("export FOO=", "Clear the FOO variable.")
			unset := r.Bool(Flag{Short: 'n', Desc: "Remove variables from the environment."})
			if err := r.Parse(); err != nil {
				return err
			}

			if *unset {
				for _, arg := range r.Args() {
					name, _ := ansi.Cut(arg, '=')
					r.sh.env.Unset(name)
				}
				return nil
			}

			for _, arg := range r.Args() {
				if !strings.Contains(arg, "=") {
					val, ok := r.sh.env.Get(arg)
					if ok {
						r.sh.env.Set(arg, val)
					}
					continue
				}

				r.sh.env.Set(ansi.Cut(arg, '='))
			}

			return nil
		}},
		Command{Name: "env", Desc: "Print shell environment values.", Exec: func(r RunArgs) error {
			if err := r.Parse(); err != nil {
				return err
			}

			for _, key := range r.sh.env.List() {
				val, _ := r.sh.env.Get(key)
				r.Println(key + "=" + val)
			}

			return nil
		}},

		Command{Name: "exit", Desc: "Exits the current shell.", Exec: func(r RunArgs) error {
			errMsg := r.String(Flag{Short: 'm', Desc: "Exit shell with an error message."})
			if err := r.Parse(); err != nil {
				return err
			}

			if errMsg != nil {
				return exitErr{errors.New(*errMsg)}
			}

			return exitErr{nil}
		}},
	)
}
