package main

import (
	"flag"
	"log"
	"os"

	"github.com/mastercactapus/embedded/bustool"
	"github.com/mastercactapus/embedded/driver/ioexp"
	"github.com/mastercactapus/embedded/serial/i2c"
	"github.com/tarm/serial"
)

func main() {
	baud := flag.Int("b", 460800, "baud rate")
	port := flag.String("p", "/dev/ttyACM0", "port")
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	p, err := serial.OpenPort(&serial.Config{Name: *port, Baud: *baud})
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	x := ioexp.NewXIAO(p)
	x.InputPins.State = 0xff
	x.OutputState.State = 0
	if err = x.Ping(); err != nil {
		log.Fatal(err)
	}
	if err = x.Flush(); err != nil {
		log.Fatal(err)
	}

	sh := bustool.NewShell(os.Stdin, os.Stdout)

	i2cSh := bustool.AddI2C(sh, i2c.New(i2c.NewSoftController(x.Pin(1), x.Pin(0))))
	bustool.AddIO(i2cSh)
	bustool.AddMem(i2cSh)

	err = sh.Run()
	if err != nil {
		log.Fatal(err)
	}
}
