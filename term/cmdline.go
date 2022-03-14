package term

import (
	"bytes"

	"github.com/mastercactapus/embedded/term/ansi"
)

// CmdLine represents a command line invocation.
type CmdLine struct {
	Args []string
	Env  []string
}

// ParseCmdLine parses a command line string.
func ParseCmdLine(cmdline string) (*CmdLine, error) {
	var p cmdParse
	p.state = cmdParam
	for _, b := range []byte(cmdline) {
		p.state = p.state(&p, b)
		if p.err != nil {
			return nil, p.err
		}
	}
	p.endOfParam()

	return &p.CmdLine, nil
}

type cmdParse struct {
	CmdLine
	state cmdParseFunc
	err   error
	buf   bytes.Buffer
	isEnv bool
	hex   byte
	hexN  int
}

type cmdParseFunc func(*cmdParse, byte) cmdParseFunc

func skipWhitespace(p *cmdParse, c byte) cmdParseFunc {
	switch c {
	case ' ', '\t', '\n', '\r':
		return skipWhitespace
	default:
		return cmdParam(p, c)
	}
}

func (p *cmdParse) endOfParam() {
	if p.buf.Len() == 0 {
		return
	}

	if p.isEnv {
		p.Env = append(p.Env, p.buf.String())
	} else {
		p.Args = append(p.Args, p.buf.String())
	}

	p.isEnv = false
	p.buf.Reset()
}

func cmdParam(p *cmdParse, c byte) cmdParseFunc {
	switch c {
	case ' ', '\t', '\n', '\r':
		p.endOfParam()
		return skipWhitespace
	case '"':
		return cmdString
	case '`':
		return cmdBackString
	case '=':
		if len(p.Args) == 0 {
			p.isEnv = true
		}
	}

	if c < ' ' || c > '~' {
		p.err = ansi.Errorf("invalid character '%c'", c)
		return nil
	}
	p.buf.WriteByte(c)
	return cmdParam
}

func cmdBackString(p *cmdParse, c byte) cmdParseFunc {
	if c == '`' {
		return cmdParam
	}

	p.buf.WriteByte(c)
	return cmdBackString
}

func cmdString(p *cmdParse, c byte) cmdParseFunc {
	switch c {
	case '\\':
		return cmdStringEsc
	case '"':
		return cmdParam
	}

	p.buf.WriteByte(c)
	return cmdString
}

func cmdStringEsc(p *cmdParse, c byte) cmdParseFunc {
	switch c {
	case '\\', '"', '\'', '`':
		p.buf.WriteByte(c)
	case 't':
		p.buf.WriteByte('\t')
	case 'n':
		p.buf.WriteByte('\n')
	case 'r':
		p.buf.WriteByte('\r')
	case 'x':
		p.hexN = 2
		return cmdStringHex
	default:
		p.err = ansi.Errorf("invalid escape sequence '\\%c'", c)
		return nil
	}

	return cmdString
}

func cmdStringHex(p *cmdParse, c byte) cmdParseFunc {
	switch {
	case '0' <= c && c <= '9':
		p.hex = p.hex*16 + (c - '0')
	case 'a' <= c && c <= 'f':
		p.hex = p.hex*16 + (c - 'a') + 10
	case 'A' <= c && c <= 'F':
		p.hex = p.hex*16 + (c - 'A') + 10
	default:
		p.err = ansi.Errorf("invalid hexadecimal character '%c'", c)
		return nil
	}

	p.hexN--
	if p.hexN == 0 {
		p.buf.WriteByte(p.hex)
		return cmdString
	}

	return cmdStringHex
}
