package term

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mastercactapus/embedded/term/ansi"
)

var ErrNotSet = fmt.Errorf("required but not set")

func (env *CommandEnv) TakeFlag(f Flag) (string, error) {
	fname := "-" + f.Name
	prefix := fname + "="

	for i, a := range env.Flags {
		switch {
		case strings.HasPrefix(a, prefix):
			_, value := ansi.Cut(a, '=')
			env.Flags = append(env.Flags[:i], env.Flags[i+1:]...)
			return value, nil
		case a == fname:
			env.Flags = append(env.Flags[:i], env.Flags[i+1:]...)
			return "", nil
		}
	}

	if f.Env != "" {
		var value string
		prefix = f.Env + "="
		for _, v := range env.LocalEnv {
			if strings.HasPrefix(v, prefix) {
				_, value = ansi.Cut(v, '=')
				return value, nil
			}
		}
		for _, v := range env.GlobalEnv {
			if strings.HasPrefix(v, prefix) {
				_, value = ansi.Cut(v, '=')
				return value, nil
			}
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

func (cmd *CmdCtx) FlagString(f Flag) string { value, _ := cmd.addFlag(f, "string"); return value }

func (cmd *CmdCtx) FlagBool(f Flag) bool {
	value, ok := cmd.addFlag(f, "bool")
	if !ok {
		return false
	}
	if value == "" {
		return true
	}
	b, err := strconv.ParseBool(value)
	cmd.setParseErr(f.valueError(err))
	return b
}

func (cmd *CmdCtx) FlagByte(f Flag) byte { return cmd.FlagUint8(f) }

func (cmd *CmdCtx) FlagInt(f Flag) int {
	value, ok := cmd.addFlag(f, "int")
	if !ok {
		return 0
	}
	i, err := strconv.Atoi(value)
	cmd.setParseErr(f.valueError(err))
	return i
}

func (cmd *CmdCtx) FlagInt8(f Flag) int8 {
	value, ok := cmd.addFlag(f, "int8")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 8)
	cmd.setParseErr(f.valueError(err))
	return int8(i)
}

func (cmd *CmdCtx) FlagInt16(f Flag) int16 {
	value, ok := cmd.addFlag(f, "int16")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 16)
	cmd.setParseErr(f.valueError(err))
	return int16(i)
}

func (cmd *CmdCtx) FlagInt32(f Flag) int32 {
	value, ok := cmd.addFlag(f, "int32")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 32)
	cmd.setParseErr(f.valueError(err))
	return int32(i)
}

func (cmd *CmdCtx) FlagInt64(f Flag) int64 {
	value, ok := cmd.addFlag(f, "int64")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 64)
	cmd.setParseErr(f.valueError(err))
	return i
}

func (cmd *CmdCtx) FlagUint(f Flag) uint {
	value, ok := cmd.addFlag(f, "uint")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 0)
	cmd.setParseErr(f.valueError(err))
	return uint(i)
}

func (cmd *CmdCtx) FlagUint8(f Flag) uint8 {
	value, ok := cmd.addFlag(f, "uint8")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 8)
	cmd.setParseErr(f.valueError(err))
	return uint8(i)
}

func (cmd *CmdCtx) FlagUint16(f Flag) uint16 {
	value, ok := cmd.addFlag(f, "uint16")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 16)
	cmd.setParseErr(f.valueError(err))
	return uint16(i)
}

func (cmd *CmdCtx) FlagUint32(f Flag) uint32 {
	value, ok := cmd.addFlag(f, "uint32")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 32)
	cmd.setParseErr(f.valueError(err))
	return uint32(i)
}

func (cmd *CmdCtx) FlagUint64(f Flag) uint64 {
	value, ok := cmd.addFlag(f, "uint64")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 64)
	cmd.setParseErr(f.valueError(err))
	return i
}
