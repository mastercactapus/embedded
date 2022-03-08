package main

import (
	"encoding/hex"
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/mastercactapus/embedded/term"
)

func main() {
	w := term.NewThrottleWriter(os.Stdout, 100)

	sh := term.NewRootShell("testshell", "Testing basic shell functionality", os.Stdin, w)
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

	subSh.AddCommand("watch", "test watch functionality", func(r term.RunArgs) error {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		rCh := r.Input()
		for {
			select {
			case <-t.C:
			case r := <-rCh:
				if r == term.Interrupt {
					return nil
				}
				continue
			}

			data := make([]byte, 64)
			rand.Read(data)
			r.Esc('H')
			r.Esc('J')
			r.Print(hex.Dump(data))
		}
	})

	err := sh.Run()
	if err != nil {
		log.Fatal(err)
	}
}
