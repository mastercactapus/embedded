package term

type Command struct {
	Name, Desc string

	Exec CmdFunc
}

type cmdData struct {
	Command
	sh      *Shell
	isShell bool
}
