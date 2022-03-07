package main

import (
	"log"
	"os"

	"github.com/mastercactapus/embedded/bus/i2c"
	"github.com/mastercactapus/embedded/bustool"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	sh := bustool.NewShell(os.Stdin, os.Stdout)
	bus := i2c.New()
	bus.Configure(i2c.Config{SDA: nilPin(false), SCL: nilPin(false)})
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
