package term

type Command struct {
	Name, Desc string

	Exec func(RunArgs) error

	sh *Shell
}
