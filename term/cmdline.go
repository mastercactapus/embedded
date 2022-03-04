package term

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// CmdLine represents a command line invocation.
type CmdLine struct {
	Args []string
	Env  []string
}

// ParseCmdLine parses a command line string.
func ParseCmdLine(cmdline string) (*CmdLine, error) {
	var cmd CmdLine
	sr := strings.NewReader(cmdline)
	for {
		arg, err := nextParam(sr)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(cmd.Args) == 0 && strings.ContainsRune(arg, '=') {
			cmd.Env = append(cmd.Env, arg)
			continue
		}

		if len(arg) == 0 && len(cmd.Args) == 0 {
			continue
		}

		cmd.Args = append(cmd.Args, arg)
	}

	return &cmd, nil
}

func nextParam(sr *strings.Reader) (string, error) {
	// skip whitespace
	for {
		r, _, err := sr.ReadRune()
		if err != nil {
			return "", err
		}
		switch r {
		case ' ', '\t', '\r', '\n':
			continue
		}

		sr.UnreadRune()
		break
	}

	var b strings.Builder
	for {
		r, _, err := sr.ReadRune()
		if err == io.EOF && b.Len() > 0 {
			return b.String(), nil
		}
		if err != nil {
			return b.String(), err
		}

		switch r {
		case '\\':
			r, _, err = sr.ReadRune()
			if err == io.EOF {
				return b.String(), io.ErrUnexpectedEOF
			}
			if err != nil {
				return b.String(), err
			}
			val, _, _, err := strconv.UnquoteChar("\\"+string(r), '\'')
			if err != nil {
				return b.String(), err
			}
			b.WriteRune(val)
		case ' ', '\t', '\n', '\r':
			return b.String(), nil
		case '"', '\'', '`':
			s, err := nextString(sr, r)
			if err != nil {
				return b.String(), err
			}
			b.WriteString(s)
		default:
			b.WriteRune(r)
		}
	}
}

func nextString(sr *strings.Reader, quote rune) (string, error) {
	var b strings.Builder
	b.WriteRune(quote)
	var escaped bool
	for {
		r, _, err := sr.ReadRune()
		if err == io.EOF {
			return "", io.ErrUnexpectedEOF
		}
		if err != nil {
			return "", err
		}

		if escaped {
			escaped = false
			b.WriteRune(r)
			continue
		}

		switch {
		case r == '\\' && quote != '`':
			escaped = true
		case r == quote:
			b.WriteRune(r)
			s, err := strconv.Unquote(b.String())
			if err != nil {
				return "", fmt.Errorf("invalid string '%s': %w", b.String(), err)
			}
			return s, nil
		}
	}
}
