package term

import (
	"fmt"

	"github.com/mastercactapus/embedded/term/ansi"
)

type RunArgs struct {
	*Flags
	*ansi.Printer

	sh *Shell
}

func (r *RunArgs) Get(k string) interface{} {
	val := r.sh.getValue(k)
	if val == nil {
		panic(fmt.Sprintf("shell=%s: cmd=%s: get '%s': not set in this or parent shell", r.sh.path(), r.Flags.cmd.Args[0], k))
	}
	return val
}
func (r *RunArgs) Set(k string, v interface{}) { r.sh.setValue(k, v) }
