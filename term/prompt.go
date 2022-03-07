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
	buf         []rune
	pos         int
	curPos      int
	lastCommand string
	hist        bool

	draw *Interruptable

	// draw interrupted
	drawInt  bool
	skipDraw bool

	stopCh chan struct{}
	drawCh chan struct{}
}

func NewPrompt(w *ansi.Printer, prompt string) *Prompt {
	p := &Prompt{
		w:      w,
		prefix: prompt,
		parse:  &ansi.Parser{},
	}
	p.draw = NewInterruptable(p.updateLine)

	return p
}

func (p *Prompt) Draw() {
	p.w.Reset()
	p.w.Print(p.prefix)
	p.w.Print(string(p.input))
	if p.pos < len(p.input) {
		p.w.CursorBack(len(p.input) - p.pos)
	}
	p.curPos = p.pos
	p.buf = append(p.buf[:0], p.input...)
}

func (p *Prompt) moveTo(i int) {
	switch {
	case i < p.curPos:
		p.w.CursorBack(p.curPos - i)
	case i > p.curPos:
		p.w.CursorForward(i - p.curPos)
	}
	p.curPos = i
}

func (p *Prompt) writeRune(r rune) {
	if p.curPos == len(p.buf) {
		p.buf = append(p.buf, r)
	} else {
		p.buf[p.curPos] = r
	}
	p.curPos++
	p.w.Print(string(r))
}

func (p *Prompt) updateLine(abort func() bool) (aborted bool) {
	if p.skipDraw {
		p.skipDraw = false
		return false
	}

	if len(p.buf) < len(p.input) {
		p.buf = append(p.buf, make([]rune, len(p.input)-len(p.buf))...)
	}

	for i := range p.input {
		if i == len(p.buf) {
			p.buf = append(p.buf, 0)
		}
		if p.buf[i] == p.input[i] {
			continue
		}
		// do stop check after so we don't mark interrupted
		// if everything is up-to-date
		if abort() {
			aborted = true
			break
		}

		p.moveTo(i)
		p.writeRune(p.input[i])
	}

	if len(p.buf) > len(p.input) {
		p.moveTo(len(p.input))
		p.w.EraseLine(ansi.CurToEnd)
		p.buf = p.buf[:len(p.input)]
	}

	p.moveTo(p.pos)

	return aborted
}

func (p *Prompt) NextCommand(r rune) (*CmdLine, error) {
	drawInt := p.draw.Interrupt()
	defer p.draw.Run()

	switch p.parse.Next(r) {
	case ansi.ValueInput:
		if p.pos == len(p.input) {
			p.input = append(p.input, r)
		} else {
			p.input = append(p.input[:p.pos], append([]rune{r}, p.input[p.pos:]...)...)
		}
		p.writeRune(r)
		p.pos++
	case ansi.ValueNewline:
		p.w.Println() // print newline first to keep it feeling responsive
		if drawInt {
			// draw sync before returning
			p.w.CursorUp(1)
			p.draw.RunSync()
			p.w.CursorDown(1)
		}
		p.hist = false
		s := strings.TrimSpace(string(p.input))
		p.input = p.input[:0]
		p.buf = p.buf[:0]
		p.pos = 0
		p.curPos = 0
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
		p.skipDraw = true
		return c, err
	case ansi.ValueCurUp:
		if p.hist || p.lastCommand == "" {
			break
		}
		p.hist = true
		p.input = append(p.input[:0], []rune(p.lastCommand)...)
		p.pos = len(p.input)
		p.moveTo(p.pos)
	case ansi.ValueCurDown:
		if !p.hist {
			break
		}
		p.hist = false
		p.input = p.input[:0]
		p.pos = 0
		p.moveTo(p.pos)
	case ansi.ValueCurLeft:
		if p.pos == 0 {
			break
		}
		p.pos--
		p.moveTo(p.pos)
	case ansi.ValueCurRight:
		if p.pos == len(p.input) {
			break
		}
		p.pos++
		p.moveTo(p.pos)
	case ansi.ValueDelPrev:
		if p.pos == 0 {
			break
		}
		p.input = append(p.input[:p.pos-1], p.input[p.pos:]...)
		p.w.CursorBack(1)
		p.curPos--
		p.writeRune(' ')
		p.w.CursorBack(1)
		p.curPos--
		p.pos--
		p.moveTo(p.pos)
	case ansi.ValueDelCur:
		if p.pos == len(p.input) {
			break
		}
		p.input = append(p.input[:p.pos], p.input[p.pos+1:]...)
		p.writeRune(' ')
	case ansi.ValueCtrlC:
		p.w.Bg(ansi.White)
		p.w.Fg(ansi.Black)
		p.w.Print("^C")
		p.w.Reset()
		p.w.Println()
		p.buf = p.buf[:0]
		p.pos = 0
		p.input = p.input[:0]
		p.skipDraw = true
		return nil, ErrInterrupt
	}

	return nil, nil
}
