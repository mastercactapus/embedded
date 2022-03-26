package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/mastercactapus/embedded/bustool"
	"github.com/mastercactapus/embedded/driver/ioexp"
	"github.com/mastercactapus/embedded/driver/stepper"
	"github.com/mastercactapus/embedded/serial/i2c"
	"github.com/mastercactapus/embedded/term"
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
	bustool.AddRTC(i2cSh)

	x.Pin(2).Output()
	x.Pin(3).Output()
	x.Pin(4).Output()
	x.Pin(5).Output()
	x.Pin(2).Set(false)
	x.Pin(3).Set(false)
	x.Pin(4).Set(false)
	x.Pin(5).Set(false)
	st := stepper.New4Phase(x.Pin(2), x.Pin(3), x.Pin(4), x.Pin(5))
	sh.AddCommand("step", "Step", func(ra term.RunArgs) error {
		n := ra.Int(term.Flag{Short: 'n', Def: "1", Desc: "number of steps"})
		freq := ra.Int(term.Flag{Short: 'f', Def: "100", Desc: "frequency (in hz)"})
		r := ra.Bool(term.Flag{Short: 'r', Desc: "reverse"})
		if err := ra.Parse(); err != nil {
			return err
		}
		st.Reverse = *r

		for i := 0; i < *n; i++ {
			err := st.Step()
			if err != nil {
				return err
			}
			time.Sleep(time.Second / time.Duration(*freq))
			if !ra.WaitForInterrupt() {
				break
			}
		}

		return st.Off()
	})

	err = sh.Run()
	if err != nil {
		log.Fatal(err)
	}
}
