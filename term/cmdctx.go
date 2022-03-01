package term

import (
	"fmt"

	"github.com/mastercactapus/embedded/term/ansi"
)

type CmdFunc func(*CmdCtx) error

type CmdCtx struct {
	p *ansi.Printer
	c *cmdData

	examples []cmdExample
	flagInfo []flagInfo
	argInfo  []argInfo

	parseErr error

	env *CommandEnv
}

func (cmd *CmdCtx) setParseErr(err error) bool {
	if cmd.parseErr != nil {
		return true
	}
	cmd.parseErr = err
	return err != nil
}

func (cmd *CmdCtx) addArgSlice(a Arg, typeName string) ([]string, bool) {
	if len(cmd.argInfo) > 0 && cmd.argInfo[len(cmd.argInfo)-1].isSlice {
		panic("cannot register arguments after a slice argument")
	}

	if a.Req && len(cmd.argInfo) > 0 && !cmd.argInfo[len(cmd.argInfo)-1].Arg.Req {
		panic("cannot register a required argument after a non-required argument")
	}

	cmd.argInfo = append(cmd.argInfo, argInfo{a, true, typeName})
	if cmd.parseErr != nil {
		return nil, false
	}

	if a.Req && len(cmd.env.Args) == 0 {
		cmd.parseErr = fmt.Errorf("argument '%s': %w", a.Name, ErrNotSet)
		return nil, false
	}

	value := cmd.env.Args
	cmd.env.Args = nil
	return value, true
}

func (cmd *CmdCtx) addArg(a Arg, typeName string) (value string, ok bool) {
	if len(cmd.argInfo) > 0 && cmd.argInfo[len(cmd.argInfo)-1].isSlice {
		panic("cannot register arguments after a slice argument")
	}

	if a.Req && len(cmd.argInfo) > 0 && !cmd.argInfo[len(cmd.argInfo)-1].Arg.Req {
		panic("cannot register a required argument after a non-required argument")
	}

	cmd.argInfo = append(cmd.argInfo, argInfo{a, false, typeName})
	if cmd.parseErr != nil {
		return "", false
	}

	if a.Req && len(cmd.env.Args) == 0 {
		cmd.parseErr = fmt.Errorf("argument '%s': %w", a.Name, ErrNotSet)
		return "", false
	}

	if len(cmd.env.Args) == 0 {
		return "", false
	}

	value = cmd.env.Args[0]
	cmd.env.Args = cmd.env.Args[1:]
	return value, true
}

func (cmd *CmdCtx) addFlag(f Flag, typeName string) (value string, ok bool) {
	cmd.flagInfo = append(cmd.flagInfo, flagInfo{f, typeName})
	if cmd.parseErr != nil {
		return "", false
	}
	value, cmd.parseErr = cmd.env.TakeFlag(f)
	return value, cmd.parseErr == nil
}

type (
	cmdExample struct{ cmdline, details string }
)

func (cmd *CmdCtx) Example(cmdline, details string) {
	cmd.examples = append(cmd.examples, cmdExample{cmdline, details})
}

func (cmd *CmdCtx) Printer() *ansi.Printer {
	return cmd.p
}

func (cmd *CmdCtx) Set(key string, v interface{}) {
	cmd.c.sh.state[key] = v
}

func (cmd *CmdCtx) Get(key string) interface{} {
	return cmd.c.sh.state[key]
}

func (cmd *CmdCtx) Parse() {
	for _, f := range cmd.env.Flags {
		if f != "-h" {
			continue
		}

		panic(usageErr{})
	}

	if cmd.parseErr != nil {
		panic(usageErr{msg: cmd.parseErr.Error()})
	}

	if len(cmd.env.Args) > 0 {
		panic(usageErr{msg: "unexpected argument: " + cmd.env.Args[0]})
	}
	if len(cmd.env.Flags) > 0 {
		panic(usageErr{msg: "unexpected flag: " + cmd.env.Flags[0]})
	}
}

func (cmd *CmdCtx) UsageError(format string, a ...interface{}) {
	panic(usageErr{msg: fmt.Sprintf(format, a...)})
}
