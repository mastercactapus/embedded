package term

type (
	exitErr  struct{ error }
	usageErr struct {
		fs  *Flags
		err error
	}
)

func (e usageErr) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	return "usage requested"
}
