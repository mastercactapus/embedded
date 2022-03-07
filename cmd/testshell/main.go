package main

import (
	"encoding/hex"
	"errors"
	"log"
	"os"

	"github.com/mastercactapus/embedded/term"
)

func main() {
	sh := term.NewRootShell("testshell", "Testing basic shell functionality", os.Stdin, os.Stdout)
	sh.AddCommand("test", "test command", func(r term.RunArgs) error {
		if err := r.Parse(); err != nil {
			return err
		}
		r.Println("hello world")
		return nil
	})
	sh.AddCommand("err", "err test command", func(r term.RunArgs) error {
		return errors.New("error")
	})

	sh.NewSubShell("subsherr", "subsh test command", func(r term.RunArgs) error {
		return errors.New("error")
	})

	subSh := sh.NewSubShell("subsh", "subsh test command", func(r term.RunArgs) error {
		r.Set("foo", "bar")
		return nil
	})

	subSh.AddCommand("test", "test command", func(r term.RunArgs) error {
		data := make([]byte, 20)
		r.Print(hex.Dump(data))
		return nil
	})
	subSh.AddCommand("err", "err test command", func(r term.RunArgs) error {
		return errors.New("error")
	})
	subSh.AddCommand("panictest", "err test command", func(r term.RunArgs) error {
		r.Get("404")
		return nil
	})

	err := sh.Run()
	if err != nil {
		log.Fatal(err)
	}
}
