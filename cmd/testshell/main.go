package main

import (
	"log"
	"os"

	"github.com/mastercactapus/embedded/term"
)

func main() {
	sh := term.NewRootShell("testshell", os.Stdin, os.Stdout)

	err := sh.Run()
	if err != nil {
		log.Fatal(err)
	}
}
