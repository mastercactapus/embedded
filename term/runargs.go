package term

import (
	"runtime"

	"github.com/mastercactapus/embedded/term/ansi"
)

type RunArgs struct {
	*Flags
	*ansi.Printer

	interrupt bool
	sh        *Shell
}

const Interrupt byte = 0x03

// Input returns a rune reader for the shell.
//
// Zero values should be ignored.
func (ra *RunArgs) Input() <-chan byte {
	return ra.sh.inputCh
}

// WaitForInterrupt will return true until CTRL+C is pressed.
func (ra *RunArgs) WaitForInterrupt() bool {
	if ra.interrupt {
		return false
	}

consumeInput:
	for {
		select {
		case c := <-ra.Input():
			if c == Interrupt {
				ra.interrupt = true
				return false
			}
			continue
		default:
			break consumeInput
		}
	}

	runtime.Gosched()
	return true
}

func (ra *RunArgs) Get(k string) interface{} {
	val := ra.sh.Get(k)
	if val == nil {
		panic("shell=" + ra.sh.path() + ": cmd=" + ra.Flags.cmd.Args[0] + ": get '" + k + "': not set in this or parent shell")
	}
	return val
}
func (ra *RunArgs) Set(k string, v interface{}) { ra.sh.Set(k, v) }
