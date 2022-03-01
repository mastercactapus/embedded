package ansi

import (
	"io"
)

type EscapeSequence struct {
	Code byte
	Args []int
}

// ParseEscapeSequence parses an escape sequence from the given reader.
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
