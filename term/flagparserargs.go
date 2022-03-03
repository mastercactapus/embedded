package term

import (
	"strconv"
	"time"
)

func (fp *FlagParser) ArgDuration(a Arg) time.Duration {
	value, ok := fp.addArg(a, "duration")
	if !ok {
		return 0
	}
	i, err := time.ParseDuration(value)
	fp.setErr(a.valueError(err))

	return i
}
func (fp *FlagParser) ArgString(a Arg) string { value, _ := fp.addArg(a, "string"); return value }
func (fp *FlagParser) ArgStringN(a Arg) []string {
	value, _ := fp.addArgSlice(a, "string")
	return value
}
func (fp *FlagParser) ArgByte(a Arg) byte { return fp.ArgUint8(a) }

func (fp *FlagParser) ArgByteN(a Arg) []byte {
	value, _ := fp.addArgSlice(a, "byte")
	result := make([]byte, len(value))
	for i, v := range value {
		b, err := strconv.ParseUint(v, 0, 8)
		if fp.setErr(a.valueError(err)) {
			break
		}
		result[i] = byte(b)
	}
	return result
}

func (fp *FlagParser) ArgBool(a Arg) bool {
	value, ok := fp.addArg(a, "bool")
	if !ok {
		return false
	}
	if value == "" {
		return false
	}

	b, err := strconv.ParseBool(value)
	fp.setErr(a.valueError(err))

	return b
}

func (fp *FlagParser) ArgInt(a Arg) int {
	value, ok := fp.addArg(a, "int")
	if !ok {
		return 0
	}
	i, err := strconv.Atoi(value)
	fp.setErr(a.valueError(err))

	return i
}

func (fp *FlagParser) ArgInt8(a Arg) int8 {
	value, ok := fp.addArg(a, "int8")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 8)
	fp.setErr(a.valueError(err))

	return int8(i)
}

func (fp *FlagParser) ArgInt16(a Arg) int16 {
	value, ok := fp.addArg(a, "int16")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 16)
	fp.setErr(a.valueError(err))

	return int16(i)
}

func (fp *FlagParser) ArgInt32(a Arg) int32 {
	value, ok := fp.addArg(a, "int32")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 32)
	fp.setErr(a.valueError(err))

	return int32(i)
}

func (fp *FlagParser) ArgInt64(a Arg) int64 {
	value, ok := fp.addArg(a, "int64")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 64)
	fp.setErr(a.valueError(err))

	return i
}

func (fp *FlagParser) ArgUint(a Arg) uint {
	value, ok := fp.addArg(a, "uint")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 0)
	fp.setErr(a.valueError(err))

	return uint(i)
}

func (fp *FlagParser) ArgUint8(a Arg) uint8 {
	value, ok := fp.addArg(a, "uint8")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 8)
	fp.setErr(a.valueError(err))

	return uint8(i)
}

func (fp *FlagParser) ArgUint16(a Arg) uint16 {
	value, ok := fp.addArg(a, "uint16")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 16)
	fp.setErr(a.valueError(err))

	return uint16(i)
}

func (fp *FlagParser) ArgUint32(a Arg) uint32 {
	value, ok := fp.addArg(a, "uint32")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 32)
	fp.setErr(a.valueError(err))

	return uint32(i)
}

func (fp *FlagParser) ArgUint64(a Arg) uint64 {
	value, ok := fp.addArg(a, "uint64")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 64)
	fp.setErr(a.valueError(err))

	return i
}
