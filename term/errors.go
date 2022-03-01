package term

type usageErr struct {
	msg string
}

func (e usageErr) Error() string {
	return e.msg
}
