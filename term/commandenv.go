package term

import (
	"fmt"
	"strings"
)

type (
	CommandEnv struct {
		Name string

		Flags     []string
		Args      []string
		LocalEnv  []string
		GlobalEnv []string
	}

	Arg struct {
		Name string
		Desc string // Description of the argument
		Req  bool   // Mark as required
	}

	Flag struct {
		Name string
		Env  string // Env var name
		Def  string // Default value
		Desc string // Description of the flag
		Req  bool   // Mark as required
	}

	flagInfo struct {
		Flag
		typeName string
	}

	argInfo struct {
		Arg
		isSlice  bool
		typeName string
	}
)

var ErrNoCommand = fmt.Errorf("no command")

func ParseCommandEnv(input string, env []string) (*CommandEnv, error) {
	args, err := SplitArgs(input)
	if err != nil {
		return nil, err
	}

	cmd := &CommandEnv{GlobalEnv: env}
	for i, a := range args {
		if strings.ContainsRune(a, '=') {
			continue
		}

		cmd.Name = a
		cmd.LocalEnv = args[:i]
		args = args[i+1:]
		break
	}
	if cmd.Name == "" {
		return nil, ErrNoCommand
	}
	for i, a := range args {
		if a == "--" {
			cmd.Flags = args[:i]
			cmd.Args = args[i+1:]
			break
		}
		if a[0] == '-' {
			cmd.Flags = args[:i+1]
			cmd.Args = args[i+1:]
		} else {
			cmd.Args = args[i:]
			break
		}
	}

	return cmd, nil
}

func (f Flag) valueError(err error) error {
	if err != nil {
		return fmt.Errorf("invalid value for flag '-%s': %w", f.Name, err)
	}

	return nil
}

func (a Arg) valueError(err error) error {
	if err != nil {
		return fmt.Errorf("invalid value for argument '%s': %w", a.Name, err)
	}

	return nil
}
