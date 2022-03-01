package term

type usageErr struct {
	err error
}

func (e usageErr) Error() string {
	if e.err != nil {
		return e.err.Error()
	}

	return "usage requested"
}
