package ansi

import (
	"fmt"
	"io"
)

type Printer struct {
	w   io.Writer
	err error
}

func NewPrinter(w io.Writer) *Printer {
	return &Printer{w: w}
}

func (pr *Printer) Indent(prefix string) *Printer {
	return NewPrinter(NewIndentWriter(pr, prefix))
}

func (pr *Printer) Err() error { return pr.err }
func (pr *Printer) Write(p []byte) (int, error) {
	if pr.err != nil {
		return 0, pr.err
	}

	n, err := pr.w.Write(p)
	pr.err = err
	return n, err
}
func (p *Printer) WriteString(s string) (int, error) { return p.Write([]byte(s)) }

func (p *Printer) Print(a ...interface{})                 { fmt.Fprint(p, a...) }
func (p *Printer) Println(a ...interface{})               { fmt.Fprintln(p, a...) }
func (p *Printer) Printf(format string, a ...interface{}) { fmt.Fprintf(p, format, a...) }

// Esc will write an escape sequence to the writer.
func (p *Printer) Esc(code byte, args ...byte) {
	p.WriteString("\x1b[")
	if len(args) == 0 || (len(args) == 1 && args[0] == 0) {
		p.Write([]byte{code})
		return
	}

	for i, arg := range args {
		if i > 0 {
			p.Write([]byte{';'})
		}
		if arg >= 100 {
			p.Write([]byte{byte(arg/100) + '0', byte(arg%100/10) + '0', byte(arg%10) + '0'})
			continue
		}
		if arg >= 10 {
			p.Write([]byte{byte(arg%100/10) + '0', byte(arg%10) + '0'})
			continue
		}
		p.Write([]byte{byte(arg) + '0'})
	}
	p.Write([]byte{code})
}

func (p *Printer) CurUp(n int) {
	if n == 1 {
		n = 0
	}
	p.Esc('A', byte(n))
}

func (p *Printer) CurDn(n int) {
	if n == 1 {
		n = 0
	}
	p.Esc('B', byte(n))
}

func (p *Printer) CurLt(n int) {
	if n == 1 {
		n = 0
	}
	p.Esc('D', byte(n))
}

func (p *Printer) CurRt(n int) {
	if n == 1 {
		n = 0
	}
	p.Esc('C', byte(n))
}

func (p *Printer) SaveCursor()    { p.Esc('s') }
func (p *Printer) RestoreCursor() { p.Esc('u') }

func (p *Printer) CurPos(x, y int) { p.Esc('H', byte(y), byte(x)) }

// Reset text color and font.
func (p *Printer) Reset() { p.Esc('m') }

// EraseLine clears the current line.
//
// If n is 0, the line is cleared from the cursor to the end of the line.
// If n is 1, the line is cleared from the cursor to the beginning of the line.
// If n is 2, the entire line is cleared.
func (p *Printer) EraseLine(n EraseLineParam) { p.Esc('K', byte(n)) }

type EraseLineParam byte

const (
	CurToEnd   EraseLineParam = 0
	CurToStart EraseLineParam = 1
	CurToAll   EraseLineParam = 2
)

// ClearLine clears the current line.
func (p *Printer) ClearLine() {
	p.WriteString("\r")
	p.Esc('K')
}

func (p *Printer) Font(f Font) {
	p.Esc('m', byte(f))
}

// Fg sets the foreground color. Use the constant values for the original 3-bit colors
// or pass in a higher value from the 8-bit color table.
func (p *Printer) Fg(c Color) {
	switch {
	case c <= White:
		p.Esc('m', 30+byte(c))
	case c > White && c <= BrightWhite:
		p.Esc('m', 90+byte(c)-8)
	default:
		p.Esc('m', 38, 5, byte(c))
	}
}

// Bg sets the background color. Use the constant values for the original 3-bit colors
// or pass in a higher value from the 8-bit color table.
func (p *Printer) Bg(c Color) {
	switch {
	case c <= White:
		p.Esc('m', 40+byte(c))
	case c > White && c <= BrightWhite:
		p.Esc('m', 100+byte(c)-8)
	default:
		p.Esc('m', 48, 5, byte(c))
	}
}

// FgRGB sets the foreground color.
func (p *Printer) FgRGB(r, g, b byte) {
	p.Esc('m', 38, 2, byte(r), byte(g), byte(b))
}

// BgRGB sets the background color.
func (p *Printer) BgRGB(r, g, b byte) {
	p.Esc('m', 38, 2, byte(r), byte(g), byte(b))
}

type Color byte

const (
	// Original 3-bit colors
	Black = Color(iota)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
)

type Font byte

const (
	Normal = Font(iota)
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	Reverse
	Conceal
	CrossedOut
)
