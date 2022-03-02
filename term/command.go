package term

import (
	"context"
	"fmt"
)

type Command struct {
	Name, Desc string

	Exec CmdFunc
	Init InitFunc
}

type (
	CmdFunc  func(context.Context) error
	InitFunc func(ctx context.Context, exec CmdFunc) error
)

func UsageError(format string, a ...interface{}) error {
	return usageErr{err: fmt.Errorf(format, a...)}
}

type cmdData struct {
	Command
	sh      *Shell
	isShell bool
}
