package term

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mastercactapus/embedded/term/ansi"
)

var ErrNotSet = fmt.Errorf("required but not set")

func (fp *FlagParser) takeFlag(f Flag) (string, error) {
	fname := "-" + f.Name
	prefix := fname + "="

	for i, a := range fp.cmd.Flags {
		switch {
		case strings.HasPrefix(a, prefix):
			_, value := ansi.Cut(a, '=')
			fp.cmd.Flags = append(fp.cmd.Flags[:i], fp.cmd.Flags[i+1:]...)
			return value, nil
		case a == fname:
			fp.cmd.Flags = append(fp.cmd.Flags[:i], fp.cmd.Flags[i+1:]...)
			return "", nil
		}
	}

	if f.Env != "" {
		var value string
		prefix = f.Env + "="
		for _, v := range fp.cmd.LocalEnv {
			if strings.HasPrefix(v, prefix) {
				_, value = ansi.Cut(v, '=')
				return value, nil
			}
		}
		if fp.lookupEnv != nil {
			value = fp.lookupEnv(f.Env)
			if value != "" {
				return value, nil
			}
		}

		if f.Def != "" {
			return f.Def, nil
		}

		if f.Req {
			return "", fmt.Errorf("flag '%s': %w", fname, ErrNotSet)
		}

		return "", nil
	}

	if f.Def != "" {
		return f.Def, nil
	}

	if f.Req {
		return "", fmt.Errorf("flag '%s': %w", fname, ErrNotSet)
	}

	return "", nil
}

func (fp *FlagParser) FlagString(f Flag) string { value, _ := fp.addFlag(f, "string"); return value }

func (fp *FlagParser) FlagBool(f Flag) bool {
	value, ok := fp.addFlag(f, "bool")
	if !ok {
		return false
	}
	if value == "" {
		return true
	}
	b, err := strconv.ParseBool(value)
	fp.setErr(f.valueError(err))
	return b
}

func (fp *FlagParser) FlagByte(f Flag) byte { return fp.FlagUint8(f) }

func (fp *FlagParser) FlagInt(f Flag) int {
	value, ok := fp.addFlag(f, "int")
	if !ok {
		return 0
	}
	i, err := strconv.Atoi(value)
	fp.setErr(f.valueError(err))
	return i
}

func (fp *FlagParser) FlagInt8(f Flag) int8 {
	value, ok := fp.addFlag(f, "int8")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 8)
	fp.setErr(f.valueError(err))
	return int8(i)
}

func (fp *FlagParser) FlagInt16(f Flag) int16 {
	value, ok := fp.addFlag(f, "int16")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 16)
	fp.setErr(f.valueError(err))
	return int16(i)
}

func (fp *FlagParser) FlagInt32(f Flag) int32 {
	value, ok := fp.addFlag(f, "int32")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 32)
	fp.setErr(f.valueError(err))
	return int32(i)
}

func (fp *FlagParser) FlagInt64(f Flag) int64 {
	value, ok := fp.addFlag(f, "int64")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 64)
	fp.setErr(f.valueError(err))
	return i
}

func (fp *FlagParser) FlagUint(f Flag) uint {
	value, ok := fp.addFlag(f, "uint")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 0)
	fp.setErr(f.valueError(err))
	return uint(i)
}

func (fp *FlagParser) FlagUint8(f Flag) uint8 {
	value, ok := fp.addFlag(f, "uint8")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 8)
	fp.setErr(f.valueError(err))
	return uint8(i)
}

func (fp *FlagParser) FlagUint16(f Flag) uint16 {
	value, ok := fp.addFlag(f, "uint16")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 16)
	fp.setErr(f.valueError(err))
	return uint16(i)
}

func (fp *FlagParser) FlagUint32(f Flag) uint32 {
	value, ok := fp.addFlag(f, "uint32")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 32)
	fp.setErr(f.valueError(err))
	return uint32(i)
}

func (fp *FlagParser) FlagUint64(f Flag) uint64 {
	value, ok := fp.addFlag(f, "uint64")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 64)
	fp.setErr(f.valueError(err))
	return i
}
