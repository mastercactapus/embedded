package term

import (
	"io"
	"strconv"
	"strings"
)

type EscapeSequence struct {
	Code byte
	Args []int
}

func (e EscapeSequence) String() string {
	var b strings.Builder
	b.WriteString("\x1b[")
	for i, arg := range e.Args {
		if i > 0 {
			b.WriteRune(';')
		}
		b.WriteString(strconv.Itoa(arg))
	}

	b.WriteByte(e.Code)
	return b.String()
}

func ParseEscapeSequence(r io.ByteReader) (*EscapeSequence, error) {
	var seq EscapeSequence
	args := make([]int, 1, 2)
	var b byte

	code, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	if code != '[' {
		return nil, nil
	}
	for {
		b, err = r.ReadByte()
		if err != nil {
			return nil, err
		}
		switch {
		case b == ';':
			args = append(args, 0)
		case b >= '0' && b <= '9':
			args[len(args)-1] = args[len(args)-1]*10 + int(b-'0')
		case b == ' ':
			continue
		default:
			seq.Code = b
			seq.Args = args
			return &seq, nil
		}
	}
}
