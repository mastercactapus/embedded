package term

import (
	"fmt"
	"runtime/debug"
)

type Command struct {
	Name, Desc string

	Exec CmdFunc
}

type cmdData struct {
	Command
	sh      *Shell
	isShell bool

	panicErr error
}

func (cmd *cmdData) exec(c *CmdCtx) error {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				cmd.panicErr = fmt.Errorf("panic: %w\n\n%s", err, debug.Stack())
			} else {
				cmd.panicErr = fmt.Errorf("panic: %v\n\n%s", r, debug.Stack())
			}
		}
	}()

	return cmd.Exec(c)
}
