package ansi

import (
	"fmt"
	"io"
)

type Printer struct {
	w   io.Writer
	err error

	lastByte byte
}

func NewPrinter(w io.Writer) *Printer {
	if ap, ok := w.(*Printer); ok {
		return ap
	}

	return &Printer{w: w}
}

func (pr *Printer) Indent(prefix string) *Printer {
	return NewPrinter(NewIndentWriter(pr, prefix))
}

func (pr *Printer) Err() error {
	err := pr.err
	pr.err = nil
	return err
}

func (pr *Printer) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if pr.err != nil {
		return 0, pr.err
	}

	var last byte = pr.lastByte
	buf := p
	for i, b := range p {
		if b != '\n' {
			last = b
			continue
		}
		if last == '\r' {
			last = b
			continue
		}

		// current is newline, but missing return
		buf = p[i:]
		w, err := pr.w.Write(p[:i])
		n += w
		if n > 0 {
			pr.lastByte = p[n-1]
		}
		if err != nil {
			pr.err = err
			return n, err
		}
		_, err = pr.w.Write([]byte{'\r'})
		if err != nil {
			pr.err = err
			return n, err
		}
		pr.lastByte = '\r'
	}

	w, err := pr.w.Write(buf)
	n += w
	pr.lastByte = buf[w-1]

	pr.err = err
	return n, err
}

func (p *Printer) WriteString(s string) (int, error) { return p.Write([]byte(s)) }

func (p *Printer) Print(a ...interface{})                 { fmt.Fprint(p, a...) }
func (p *Printer) Println(a ...interface{})               { fmt.Fprintln(p, a...) }
func (p *Printer) Printf(format string, a ...interface{}) { fmt.Fprintf(p, format, a...) }

// Esc will write an escape sequence to the writer.
func (p *Printer) Esc(code byte, args ...int) {
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

func (p *Printer) CursorUp(n int) {
	if n == 1 {
		n = 0
	}
	p.Esc('A', n)
}

func (p *Printer) CursorDown(n int) {
	if n == 1 {
		n = 0
	}
	p.Esc('B', n)
}

func (p *Printer) CursorForward(n int) {
	if n == 1 {
		n = 0
	}
	p.Esc('C', n)
}

func (p *Printer) CursorBack(n int) {
	if n == 1 {
		p.Print("\b")
		return
	}
	p.Esc('D', n)
}

func (p *Printer) SaveCursor()    { p.Esc('7') }
func (p *Printer) RestoreCursor() { p.Esc('8') }

func (p *Printer) CurPos(x, y int) { p.Esc('H', y, x) }

// Reset text color and font.
func (p *Printer) Reset() { p.Esc('m') }

// EraseLine clears the current line.
//
// If n is 0, the line is cleared from the cursor to the end of the line.
// If n is 1, the line is cleared from the cursor to the beginning of the line.
// If n is 2, the entire line is cleared.
func (p *Printer) EraseLine(n EraseLineParam) { p.Esc('K', int(n)) }

type EraseLineParam int

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
	p.Esc('m', int(f))
}

// Fg sets the foreground color. Use the constant values for the original 3-bit colors
// or pass in a higher value from the 8-bit color table.
func (p *Printer) Fg(c Color) {
	switch {
	case c <= White:
		p.Esc('m', 30+int(c))
	case c > White && c <= BrightWhite:
		p.Esc('m', 90+int(c)-8)
	default:
		p.Esc('m', 38, 5, int(c))
	}
}

// Bg sets the background color. Use the constant values for the original 3-bit colors
// or pass in a higher value from the 8-bit color table.
func (p *Printer) Bg(c Color) {
	switch {
	case c <= White:
		p.Esc('m', int(40+c))
	case c > White && c <= BrightWhite:
		p.Esc('m', int(100+c-8))
	default:
		p.Esc('m', 48, 5, int(c))
	}
}

// FgRGB sets the foreground color.
func (p *Printer) FgRGB(r, g, b uint8) {
	p.Esc('m', 38, 2, int(r), int(g), int(b))
}

// BgRGB sets the background color.
func (p *Printer) BgRGB(r, g, b uint8) {
	p.Esc('m', 38, 2, int(r), int(g), int(b))
}

type Color uint8

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

type Font int

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
