package term

import "fmt"

func NewFlagParser(cmd *CommandEnv, lookupEnv func(string) string) *FlagParser {
	fp := &FlagParser{
		cmd: cmd,

		lookupEnv: lookupEnv,
	}

	fp.showHelp = fp.FlagBool(Flag{Name: "h", Desc: "Show this help message"})
	return fp
}

type cmdExample struct {
	cmdline string
	details string
}

type FlagParser struct {
	cmd *CommandEnv

	examples []cmdExample
	flagInfo []flagInfo
	argInfo  []argInfo

	lookupEnv func(string) string
	err       error

	showHelp bool
}

func (fp *FlagParser) Err() error {
	if fp.showHelp {
		return usageErr{fp, nil}
	}

	if fp.err != nil {
		return usageErr{fp, fp.err}
	}

	if len(fp.cmd.Args) > 0 {
		return usageErr{fp, fmt.Errorf("unexpected argument: %s", fp.cmd.Args[0])}
	}
	if len(fp.cmd.Flags) > 0 {
		return usageErr{fp, fmt.Errorf("unexpected flag: %s", fp.cmd.Flags[0])}
	}

	return nil
}

func (fp *FlagParser) UsageError(format string, args ...interface{}) error {
	return usageErr{fp, fmt.Errorf(format, args...)}
}

func (fp *FlagParser) Example(cmdline, details string) {
	fp.examples = append(fp.examples, cmdExample{cmdline, details})
}

func (fp *FlagParser) setErr(err error) bool {
	if fp.err != nil {
		return true
	}

	fp.err = err
	return err != nil
}

func (fp *FlagParser) addArgSlice(a Arg, typeName string) ([]string, bool) {
	if len(fp.argInfo) > 0 && fp.argInfo[len(fp.argInfo)-1].isSlice {
		panic("cannot register arguments after a slice argument")
	}

	if a.Req && len(fp.argInfo) > 0 && !fp.argInfo[len(fp.argInfo)-1].Arg.Req {
		panic("cannot register a required argument after a non-required argument")
	}

	fp.argInfo = append(fp.argInfo, argInfo{a, true, typeName})
	if fp.err != nil {
		return nil, false
	}

	if a.Req && len(fp.cmd.Args) == 0 {
		fp.err = fmt.Errorf("argument '%s': %w", a.Name, ErrNotSet)
		return nil, false
	}

	value := fp.cmd.Args
	fp.cmd.Args = nil
	return value, true
}

func (fp *FlagParser) addArg(a Arg, typeName string) (value string, ok bool) {
	if len(fp.argInfo) > 0 && fp.argInfo[len(fp.argInfo)-1].isSlice {
		panic("cannot register arguments after a slice argument")
	}

	if a.Req && len(fp.argInfo) > 0 && !fp.argInfo[len(fp.argInfo)-1].Arg.Req {
		panic("cannot register a required argument after a non-required argument")
	}

	fp.argInfo = append(fp.argInfo, argInfo{a, false, typeName})
	if fp.err != nil {
		return "", false
	}

	if a.Req && len(fp.cmd.Args) == 0 {
		fp.err = fmt.Errorf("argument '%s': %w", a.Name, ErrNotSet)
		return "", false
	}

	if len(fp.cmd.Args) == 0 {
		return "", false
	}

	value = fp.cmd.Args[0]
	fp.cmd.Args = fp.cmd.Args[1:]
	return value, true
}

func (fp *FlagParser) addFlag(f Flag, typeName string) (value string, ok bool) {
	fp.flagInfo = append(fp.flagInfo, flagInfo{f, typeName})
	if fp.err != nil {
		return "", false
	}
	var isSet bool
	value, isSet, fp.err = fp.takeFlag(f)
	if typeName == "bool" && !isSet {
		return "", false
	}
	return value, fp.err == nil
}
