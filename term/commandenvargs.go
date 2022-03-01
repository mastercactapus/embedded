package term

import "strconv"

func (cmd *CmdCtx) ArgString(a Arg) string    { value, _ := cmd.addArg(a, "string"); return value }
func (cmd *CmdCtx) ArgStringN(a Arg) []string { value, _ := cmd.addArgSlice(a, "string"); return value }
func (cmd *CmdCtx) ArgByte(a Arg) byte        { return cmd.ArgUint8(a) }

func (cmd *CmdCtx) ArgByteN(a Arg) []byte {
	value, _ := cmd.addArgSlice(a, "byte")
	result := make([]byte, len(value))
	for i, v := range value {
		b, err := strconv.ParseUint(v, 0, 8)
		if cmd.setParseErr(a.valueError(err)) {
			break
		}
		result[i] = byte(b)
	}
	return result
}

func (cmd *CmdCtx) ArgBool(a Arg) bool {
	value, ok := cmd.addArg(a, "bool")
	if !ok {
		return false
	}
	if value == "" {
		return false
	}

	b, err := strconv.ParseBool(value)
	cmd.setParseErr(a.valueError(err))

	return b
}

func (cmd *CmdCtx) ArgInt(a Arg) int {
	value, ok := cmd.addArg(a, "int")
	if !ok {
		return 0
	}
	i, err := strconv.Atoi(value)
	cmd.setParseErr(a.valueError(err))

	return i
}

func (cmd *CmdCtx) ArgInt8(a Arg) int8 {
	value, ok := cmd.addArg(a, "int8")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 8)
	cmd.setParseErr(a.valueError(err))

	return int8(i)
}

func (cmd *CmdCtx) ArgInt16(a Arg) int16 {
	value, ok := cmd.addArg(a, "int16")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 16)
	cmd.setParseErr(a.valueError(err))

	return int16(i)
}

func (cmd *CmdCtx) ArgInt32(a Arg) int32 {
	value, ok := cmd.addArg(a, "int32")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 32)
	cmd.setParseErr(a.valueError(err))

	return int32(i)
}

func (cmd *CmdCtx) ArgInt64(a Arg) int64 {
	value, ok := cmd.addArg(a, "int64")
	if !ok {
		return 0
	}
	i, err := strconv.ParseInt(value, 0, 64)
	cmd.setParseErr(a.valueError(err))

	return i
}

func (cmd *CmdCtx) ArgUint(a Arg) uint {
	value, ok := cmd.addArg(a, "uint")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 0)
	cmd.setParseErr(a.valueError(err))

	return uint(i)
}

func (cmd *CmdCtx) ArgUint8(a Arg) uint8 {
	value, ok := cmd.addArg(a, "uint8")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 8)
	cmd.setParseErr(a.valueError(err))

	return uint8(i)
}

func (cmd *CmdCtx) ArgUint16(a Arg) uint16 {
	value, ok := cmd.addArg(a, "uint16")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 16)
	cmd.setParseErr(a.valueError(err))

	return uint16(i)
}

func (cmd *CmdCtx) ArgUint32(a Arg) uint32 {
	value, ok := cmd.addArg(a, "uint32")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 32)
	cmd.setParseErr(a.valueError(err))

	return uint32(i)
}

func (cmd *CmdCtx) ArgUint64(a Arg) uint64 {
	value, ok := cmd.addArg(a, "uint64")
	if !ok {
		return 0
	}
	i, err := strconv.ParseUint(value, 0, 64)
	cmd.setParseErr(a.valueError(err))

	return i
}
