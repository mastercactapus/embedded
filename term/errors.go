package term

import "io"

type usageErr struct {
	fp  *FlagParser
	err error
}

func (e usageErr) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	return "usage requested"
}
func (e usageErr) PrintUsage(w io.Writer) { e.fp.PrintUsage(w) }
