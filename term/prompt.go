package term

import (
	"errors"
	"strings"

	"github.com/mastercactapus/embedded/term/ansi"
)

var ErrInterrupt = errors.New("interrupted")

type Prompt struct {
	w     *ansi.Printer
	parse *ansi.Parser

	prefix      string
	input       []rune
	pos         int
	lastCommand string
	hist        bool
}

func NewPrompt(w *ansi.Printer, prompt string) *Prompt {
	p := &Prompt{
		w:      w,
		prefix: prompt,
		parse:  &ansi.Parser{},
	}

	return p
}

func (p *Prompt) Draw() {
	p.w.Reset()
	p.w.Print(p.prefix)
	p.w.Print(string(p.input))
	if p.pos < len(p.input) {
		p.w.CursorBack(len(p.input) - p.pos)
	}
	p.w.Esc('h', 4)
}

func (p *Prompt) NextCommand(r rune) (*CmdLine, error) {
	if r == 0 {
		// special case, ignore entirely
		// used to signal we are processing input
		return nil, nil
	}

	switch p.parse.Next(r) {
	case ansi.ValueInput:
		if p.pos == len(p.input) {
			p.input = append(p.input, r)
			p.w.WriteString(string(r))
			p.pos++
			break
		}
		p.input = append(p.input[:p.pos], append([]rune{r}, p.input[p.pos:]...)...)
		p.w.WriteString(string(r))
		p.pos++
	case ansi.ValueNewline:
		p.w.Println() // print newline first to keep it feeling responsive
		p.w.Esc('l', 4)
		p.hist = false
		s := strings.TrimSpace(string(p.input))
		p.input = p.input[:0]
		p.pos = 0
		if len(s) == 0 {
			p.Draw()
			break
		}

		c, err := ParseCmdLine(s)
		if err != nil {
			return nil, err
		}

		if len(c.Args) == 0 {
			// just ENV, ignore
			p.Draw()
			break
		}

		p.lastCommand = s
		return c, err
	case ansi.ValueCurUp:
		if p.hist || p.lastCommand == "" {
			break
		}
		p.hist = true
		p.input = append(p.input[:0], []rune(p.lastCommand)...)
		if p.pos > 0 {
			p.w.CursorBack(p.pos)
		}
		p.w.Print(string(p.input))
		p.w.EraseLine(ansi.CurToEnd)
		p.pos = len(p.input)
	case ansi.ValueCurDown:
		if !p.hist {
			break
		}
		p.hist = false
		if p.pos > 0 {
			p.w.CursorBack(p.pos)
		}
		p.w.EraseLine(ansi.CurToEnd)
		p.input = p.input[:0]
		p.pos = 0
	case ansi.ValueCurLeft:
		if p.pos == 0 {
			break
		}
		p.pos--
		p.w.CursorBack(1)
	case ansi.ValueCurRight:
		if p.pos == len(p.input) {
			break
		}
		p.pos++
		p.w.CursorForward(1)
	case ansi.ValueDelPrev:
		if p.pos == 0 {
			break
		}
		p.input = append(p.input[:p.pos-1], p.input[p.pos:]...)
		p.w.WriteByte('\b')
		p.w.Esc('P')
		p.pos--
	case ansi.ValueDelCur:
		if p.pos == len(p.input) {
			break
		}
		p.input = append(p.input[:p.pos], p.input[p.pos+1:]...)
		p.w.Esc('P')
	case ansi.ValueCtrlC:
		p.w.Bg(ansi.White)
		p.w.Fg(ansi.Black)
		p.w.Print("^C")
		p.w.Reset()
		p.w.Println()
		p.pos = 0
		p.input = p.input[:0]
		return nil, ErrInterrupt
	}

	return nil, nil
}
