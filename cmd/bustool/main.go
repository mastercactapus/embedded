package main

import (
	"context"
	"log"
	"os"

	"github.com/mastercactapus/embedded/bustool"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	sh := bustool.NewShell(os.Stdin, os.Stdout)
	i2cSh := bustool.AddI2C(sh, nilPin(false), nilPin(true))
	bustool.AddMem(i2cSh)
	bustool.AddIO(i2cSh)

	s, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatalln(err)
	}
	defer terminal.Restore(0, s)

	err = sh.Exec(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
