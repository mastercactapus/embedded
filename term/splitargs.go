package term

import (
	"strconv"
)

func SplitArgs(input string) ([]string, error) {
	var args []string
	var arg string
	var quote rune
	for _, r := range input {
		if quote == 0 {
			switch r {
			case ' ', '\t', '\n', '\r':
				if len(arg) > 0 {
					if arg[0] != '"' && arg[0] != '\'' && arg[0] != '`' {
						arg = `"` + arg + `"`
					}
					a, err := strconv.Unquote(string(arg))
					if err != nil {
						return nil, err
					}

					args = append(args, a)
					arg = ""
				}
			case '"', '\'', '`':
				quote = r
				arg += string(r)
			default:
				arg += string(r)
			}
		} else {
			if r == quote {
				quote = 0
			}
			arg += string(r)
		}
	}
	if len(arg) > 0 {
		if arg[0] != '"' && arg[0] != '\'' && arg[0] != '`' {
			arg = `"` + arg + `"`
		}
		a, err := strconv.Unquote(string(arg))
		if err != nil {
			return nil, err
		}

		args = append(args, a)
	}

	return args, nil
}
