package term

import (
	"context"
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

type cmdData struct {
	Command
	sh      *Shell
	isShell bool
}
