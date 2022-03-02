package term

import (
	"bufio"
	"context"
	"io"
	"sort"

	"github.com/mastercactapus/embedded/term/ansi"
)

type Shell struct {
	parent *Shell

	name string
	desc string

	r *bufio.Reader
	w *newliner
	p *ansi.Printer

	lastWByte byte

	err error

	lastCmdErr error

	bNames   []string
	shNames  []string
	cmdNames []string

	cmds map[string]*cmdData

	env []string
}

func (sh *Shell) dir() string {
	if sh.parent == nil {
		return "/"
	}

	return sh.parent.dir() + sh.parent.name + "/"
}

func NewShell(name, desc string, r io.Reader, w io.Writer) *Shell {
	if name == "" {
		name = "default"
	}

	nl, ok := w.(*newliner)
	if !ok {
		nl = &newliner{Writer: bufio.NewWriter(w)}
	}

	sh := &Shell{
		name: name,
		desc: desc,
		r:    bufio.NewReader(r),
		w:    nl,
		cmds: make(map[string]*cmdData),
	}
	sh.p = ansi.NewPrinter(sh.w)

	for _, cmd := range builtins {
		sh.cmds[cmd.Name] = &cmdData{Command: cmd, sh: sh}
		sh.bNames = append(sh.bNames, cmd.Name)
	}

	return sh
}

type exitErr struct{ error }

func launchSubShell(ctx context.Context) error {
	if err := Parse(ctx).Err(); err != nil {
		return err
	}

	cmd := ctx.Value(ctxKeyCmd).(*cmdContext)
	return cmd.sh.Exec(ctx)
}

func (sh *Shell) NewSubShell(cmd Command) *Shell {
	if _, ok := sh.cmds[cmd.Name]; ok {
		panic("command already exists: " + cmd.Name)
	}

	subSh := NewShell(cmd.Name, cmd.Desc, sh.r, sh.w)
	subSh.parent = sh
	cmd.Exec = launchSubShell

	sh.cmds[cmd.Name] = &cmdData{Command: cmd, sh: subSh, isShell: true}
	sh.shNames = append(sh.shNames, cmd.Name)
	sort.Strings(sh.shNames)

	return subSh
}

func (sh *Shell) AddCommand(cmd Command) {
	if _, ok := sh.cmds[cmd.Name]; ok {
		panic("command already exists: " + cmd.Name)
	}

	sh.cmds[cmd.Name] = &cmdData{Command: cmd, sh: sh}

	sh.cmdNames = append(sh.cmdNames, cmd.Name)
	sort.Strings(sh.cmdNames)
}

func (sh *Shell) prompt(input string) {
	defer sh.w.Flush()

	sh.p.ClearLine()
	sh.p.Reset()
	sh.p.Print(sh.dir() + sh.name)

	if sh.lastCmdErr != nil {
		sh.p.Fg(ansi.Red)
		sh.p.Print(" [")
		sh.p.Font(ansi.Bold)
		sh.p.Print("ERR")
		sh.p.Font(ansi.Normal)
		sh.p.Fg(ansi.Red)
		sh.p.Print("]")
		sh.p.Reset()
	}

	sh.p.Print("> " + input)
}

func (sh *Shell) WriteByte(p byte) error {
	if sh.err != nil {
		return sh.err
	}
	if p == 0 {
		return nil
	}

	if p == '\n' && sh.lastWByte != '\r' {
		sh.w.WriteByte('\r')
	}

	return sh.w.WriteByte(p)
}

func (sh *Shell) Write(p []byte) (int, error) {
	if sh.err != nil {
		return 0, sh.err
	}

	for _, b := range p {
		sh.err = sh.WriteByte(b)
	}

	return len(p), sh.err
}
