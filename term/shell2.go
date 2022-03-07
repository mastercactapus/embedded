package term

import (
	"bufio"
	"errors"
	"io"

	"github.com/mastercactapus/embedded/term/ansi"
)

type Shell2 struct {
	name string

	parent *Shell2
	prompt *Prompt
	w      *ansi.Printer
	r      io.RuneReader

	commands []Command2
	shells   []Command2

	env    *Env
	values map[string]interface{}
}

func NewRootShell(name string, in io.Reader, out io.Writer) *Shell2 {
	p := ansi.NewPrinter(out)
	sh := &Shell2{
		name:   name,
		w:      p,
		r:      bufio.NewReader(in),
		prompt: NewPrompt(p, name+"> "),
		env:    NewEnv(),
	}

	return sh
}

func (sh *Shell2) Run() error {
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

func (sh *Shell2) findCommand(name string) *Command2 {
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

func (sh *Shell2) printUsage(cmd *Command2, usage usageErr) {
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

func (sh *Shell2) runCommand(cmd *Command2, cmdline *CmdLine) error {
	return cmd.Exec(RunArgs{
		Flags2:  NewFlagSet(cmdline, sh.env.Get),
		Printer: sh.w,
		sh:      sh,
	})
}

type Command2 struct {
	Name, Desc string

	Exec func(RunArgs) error

	sh *Shell2
}

type RunArgs struct {
	*Flags2
	*ansi.Printer

	sh *Shell2
}

func (r *RunArgs) Get(k string) interface{}    { return r.sh.getValue(k) }
func (r *RunArgs) Set(k string, v interface{}) { r.sh.setValue(k, v) }

func (sh *Shell2) setValue(k string, v interface{}) {
	if sh.values == nil {
		sh.values = make(map[string]interface{})
	}
	sh.values[k] = v
}

func (sh *Shell2) getValue(k string) interface{} {
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

func (sh *Shell2) path() string {
	var parentName string
	if sh.parent != nil {
		return sh.parent.path()
	}

	return parentName + "/" + sh.name
}

func (sh *Shell2) AddCommand(cmd Command2) {
	if sh.findCommand(cmd.Name) != nil {
		panic("Duplicate command: " + cmd.Name)
	}
	sh.commands = append(sh.commands, cmd)
}

func (sh *Shell2) AddSubShell(cmd Command2) {
	if sh.findCommand(cmd.Name) != nil {
		panic("Duplicate command: " + cmd.Name)
	}
	cmd.sh = &Shell2{
		parent: sh,
		w:      sh.w,
		env:    NewEnv(),
	}
	cmd.sh.env.SetParent(sh.env)
	cmd.sh.prompt = NewPrompt(sh.w, cmd.sh.path()+"> ")
	sh.shells = append(sh.shells, cmd)
}
