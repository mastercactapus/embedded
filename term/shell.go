package term

import (
	"bufio"
	"errors"
	"io"

	"github.com/mastercactapus/embedded/term/ansi"
)

type Shell struct {
	name, desc string

	parent *Shell
	prompt *Prompt
	w      *ansi.Printer
	r      io.RuneReader

	commands []Command
	shells   []Command

	noExit bool

	env    *Env
	values map[string]interface{}
}

func NewRootShell(name, desc string, in io.Reader, out io.Writer) *Shell {
	p := ansi.NewPrinter(out)
	sh := &Shell{
		name: name,
		desc: desc,
		w:    p,
		r:    bufio.NewReader(&fixReader{in}),
		env:  NewEnv(),
	}
	sh.prompt = NewPrompt(p, sh.path()+"> ")

	return sh
}

func (sh *Shell) SetNoExit(v bool) { sh.noExit = v }

func (sh *Shell) Run() error {
	sh.prompt.Draw()
	for {
		r, _, err := sh.r.ReadRune()
		if err != nil {
			return err
		}

		cmdLine, err := sh.prompt.NextCommand(r)
		if errors.Is(err, ErrInterrupt) {
			sh.prompt.Draw()
			continue
		}
		if cmdLine == nil {
			continue
		}

		cmd := sh.findCommand(cmdLine.Args[0])
		if cmd == nil {
			sh.w.Println("Unknown command: '" + cmdLine.Args[0] + "' try 'help'.")
			sh.prompt.Draw()
			continue
		}

		err = sh.runCommand(cmd, cmdLine)
		var exit exitErr
		var usage usageErr
		switch {
		case errors.As(err, &exit):
			return exit.error
		case errors.As(err, &usage):
			sh.printUsage(cmd, usage)
		case err != nil:
			sh.w.Fg(ansi.Red)
			sh.w.Println(err)
		}
		sh.prompt.Draw()
	}
}

func (sh *Shell) findCommand(name string) *Command {
	if sh.noExit && name == "exit" {
		return nil
	}

	for _, cmd := range sh.commands {
		if cmd.Name != name {
			continue
		}

		return &cmd
	}

	for _, cmd := range sh.shells {
		if cmd.Name != name {
			continue
		}

		return &cmd
	}

	for _, cmd := range builtin {
		if cmd.Name != name {
			continue
		}

		return &cmd
	}

	return nil
}

func (sh *Shell) printUsage(cmd *Command, usage usageErr) {
	if usage.err != nil {
		sh.w.Fg(ansi.Red)
		sh.w.Println(usage.err.Error())
		sh.w.Reset()
	}

	if cmd.Desc != "" {
		sh.w.Println(cmd.Desc)
		sh.w.Println()
	}

	usage.PrintUsage(sh.w)
}

func (sh *Shell) runCommand(cmd *Command, cmdline *CmdLine) error {
	if cmd.Exec != nil {
		err := cmd.Exec(RunArgs{
			Flags:   NewFlagSet(cmdline, sh.env.Get),
			Printer: sh.w,
			sh:      sh,
		})
		if err != nil {
			return err
		}
	}

	if cmd.sh == nil {
		return nil
	}

	return cmd.sh.Run()
}

func (sh *Shell) setValue(k string, v interface{}) {
	if sh.values == nil {
		sh.values = make(map[string]interface{})
	}
	sh.values[k] = v
}

func (sh *Shell) getValue(k string) interface{} {
	if sh.values != nil {
		if v, ok := sh.values[k]; ok {
			return v
		}
	}

	if sh.parent != nil {
		return sh.parent.getValue(k)
	}

	return nil
}

func (sh *Shell) path() string {
	var parentName string
	if sh.parent != nil {
		parentName = sh.parent.path()
	}

	return parentName + "/" + sh.name
}

func (sh *Shell) AddCommand(cmd Command) {
	if sh.findCommand(cmd.Name) != nil {
		panic("Duplicate command: " + cmd.Name)
	}
	sh.commands = append(sh.commands, cmd)
}

func (sh *Shell) NewSubShell(cmd Command) *Shell {
	if sh.findCommand(cmd.Name) != nil {
		panic("Duplicate command: " + cmd.Name)
	}
	cmd.sh = &Shell{
		name:   cmd.Name,
		desc:   cmd.Desc,
		parent: sh,
		w:      sh.w,
		r:      sh.r,
		env:    NewEnv(),
	}
	cmd.sh.env.SetParent(sh.env)
	cmd.sh.prompt = NewPrompt(sh.w, cmd.sh.path()+"> ")
	sh.shells = append(sh.shells, cmd)
	return cmd.sh
}
