package term

import (
	"bufio"
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

	init CmdFunc

	lastWByte byte

	err error

	lastCmdErr error

	bNames   []string
	shNames  []string
	cmdNames []string

	cmds  map[string]*cmdData
	state map[string]interface{}

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

func launchSubShell(c *CmdCtx) error {
	if err := c.Parse(); err != nil {
		return err
	}
	return c.c.sh.Exec()
}

func (sh *Shell) NewSubShell(cmd Command) *Shell {
	if _, ok := sh.cmds[cmd.Name]; ok {
		panic("command already exists: " + cmd.Name)
	}

	subSh := NewShell(cmd.Name, cmd.Desc, sh.r, sh.w)
	subSh.parent = sh
	subSh.init = cmd.Exec
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

func (s *Shell) prompt(input string) {
	defer s.w.Flush()

	s.p.ClearLine()
	s.p.Reset()
	s.p.Print(s.dir() + s.name)

	if s.lastCmdErr != nil {
		s.p.Fg(ansi.Red)
		s.p.Print(" [")
		s.p.Font(ansi.Bold)
		s.p.Print("ERR")
		s.p.Font(ansi.Normal)
		s.p.Fg(ansi.Red)
		s.p.Print("]")
		s.p.Reset()
	}

	s.p.Print("> " + input)
}

func (t *Shell) WriteByte(p byte) error {
	if t.err != nil {
		return t.err
	}
	if p == 0 {
		return nil
	}

	if p == '\n' && t.lastWByte != '\r' {
		t.w.WriteByte('\r')
	}

	return t.w.WriteByte(p)
}

func (t *Shell) Write(p []byte) (int, error) {
	if t.err != nil {
		return 0, t.err
	}

	for _, b := range p {
		t.err = t.WriteByte(b)
	}

	return len(p), t.err
}
