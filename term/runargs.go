package term

import (
	"fmt"
	"runtime"

	"github.com/mastercactapus/embedded/term/ansi"
)

type RunArgs struct {
	*Flags
	*ansi.Printer

	sh *Shell
}

const Interrupt rune = 0x03

// Input returns a rune reader for the shell.
//
// Zero values should be ignored.
func (r *RunArgs) Input() <-chan rune {
	return r.sh.r
}

// WaitForInterrupt will return true until CTRL+C is pressed.
func (r *RunArgs) WaitForInterrupt() bool {
	select {
	case c := <-r.Input():
		if c == Interrupt {
			return false
		}
	default:
	}

	runtime.Gosched()
	return true
}

func (r *RunArgs) Get(k string) interface{} {
	val := r.sh.Get(k)
	if val == nil {
		panic(fmt.Sprintf("shell=%s: cmd=%s: get '%s': not set in this or parent shell", r.sh.path(), r.Flags.cmd.Args[0], k))
	}
	return val
}
func (r *RunArgs) Set(k string, v interface{}) { r.sh.Set(k, v) }
