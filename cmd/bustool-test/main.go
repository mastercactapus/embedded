package main

import (
	"log"
	"os"

	"github.com/mastercactapus/embedded/bustool"
	"github.com/mastercactapus/embedded/serial/i2c"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	sh := bustool.NewShell(os.Stdin, os.Stdout)
	bus := i2c.New(i2c.NewSoftController(nilPin(false), nilPin(false)))
	i2cSh := bustool.AddI2C(sh, bus)
	bustool.AddMem(i2cSh)
	bustool.AddIO(i2cSh)

	s, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatalln(err)
	}
	defer terminal.Restore(0, s)

	err = sh.Run()
	if err != nil {
		log.Fatal(err)
	}
}
