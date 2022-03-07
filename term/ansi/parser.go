package ansi

import (
	"unicode"
)

type ParserValueType int

const (
	ValueNone ParserValueType = iota

	ValueInput // printable input character
	ValueCtrlC
	ValueCtrlD
	ValueNewline
	ValueDelPrev // delete previous character
	ValueDelCur  // delete current character
	ValueCurUp
	ValueCurDown
	ValueCurLeft
	ValueCurRight
)

// Parser will parse user input rune-by-rune.
type Parser struct {
	state   parserState
	typ     ParserValueType
	csiArgs []int

	lastNewline rune
}

// Args returns the parsed CSI arguments.
//
// It will be set with a length >= 1 for cursor value types
// with the first element being >= 1. The second element is
// for modifier keys.
//
// It is invalid to use the returned slice after the Next is called again.
func (p *Parser) Args() []int {
	return p.csiArgs
}

func (p *Parser) Next(r rune) ParserValueType {
	if p.state == nil {
		p.state = stateInput
	}
	p.typ = ValueNone
	p.state = p.state(p, r)
	return p.typ
}

type parserState func(*Parser, rune) parserState

func stateInput(p *Parser, r rune) parserState {
	switch r {
	case 0x03: // ctrl+c
		p.typ = ValueCtrlC
	case 0x04: // ctrl+d
		p.typ = ValueCtrlD
	case '\r', '\n':
		if p.lastNewline != r && p.lastNewline != 0 {
			p.lastNewline = r
			// ignore
			return stateInput
		}
		p.lastNewline = r
		p.typ = ValueNewline
		return stateInput
	case 0x7f, 0x08: // backspace
		p.typ = ValueDelPrev
	case 0x1b:
		return stateEsc
	}

	if unicode.IsPrint(r) {
		p.typ = ValueInput
	}

	return stateInput
}

func stateEsc(p *Parser, r rune) parserState {
	switch r {
	case '[':
		p.csiArgs = append(p.csiArgs[:0], 0)
		return stateCSI
	case 'P', ']', 'X', '^', '_':
		// ignore until ST
		return stateEscStrIgnore
	}

	return stateInput
}

func stateCSI(p *Parser, r rune) parserState {
	if r >= '0' && r <= '9' {
		p.csiArgs[len(p.csiArgs)-1] = p.csiArgs[len(p.csiArgs)-1]*10 + int(r-'0')
		return stateCSI
	}

	switch r {
	case ';':
		p.csiArgs = append(p.csiArgs, 0)
		return stateCSI
	case '~':
		p.typ = ValueDelCur
		return stateInput
	case 'A', 'B', 'C', 'D':
		return stateCursorCSI(p, r)
	}

	return stateCSIParamIgnore(p, r)
}

func stateCursorCSI(p *Parser, r rune) parserState {
	if p.csiArgs[0] == 0 {
		p.csiArgs[0] = 1
	}
	switch r {
	case 'A':
		p.typ = ValueCurUp
	case 'B':
		p.typ = ValueCurDown
	case 'C':
		p.typ = ValueCurRight
	case 'D':
		p.typ = ValueCurLeft
	}

	return stateInput
}

func stateCSIParamIgnore(p *Parser, r rune) parserState {
	if r >= '0' && r <= '?' {
		return stateCSIParamIgnore
	}

	return stateCSI
}

func stateCSIItermediateIgnore(p *Parser, r rune) parserState {
	if r >= ' ' && r <= '/' {
		return stateCSIItermediateIgnore
	}

	// invalid or done
	return stateInput
}

func stateEscStrIgnore(p *Parser, r rune) parserState {
	if r == 0x1b {
		return stateEscStrIgnoreEsc
	}

	return stateEscStrIgnore
}

func stateEscStrIgnoreEsc(p *Parser, r rune) parserState {
	if r == '\\' {
		return stateInput
	}

	return stateEscStrIgnore
}
