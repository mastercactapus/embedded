package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"machine"

	"github.com/mastercactapus/embedded/driver"
	"github.com/mastercactapus/embedded/serial/i2c"
)

type fixReader struct {
	io.Reader
}

func (r *fixReader) Read(p []byte) (n int, err error) {
	for {
		n, err = r.Reader.Read(p)
		if err != nil {
			return
		}

		// Serial reader is broken, so we need to fix it
		if n > 0 {
			return
		}
	}
}

var pins = [12]machine.Pin{
	machine.D0,
	machine.D1,
	machine.D2,
	machine.D3,
	machine.D4,
	machine.D5,
	machine.D6,
	machine.D7,
	machine.D8,
	machine.D9,
	machine.D10,
	machine.LED,
}

var i2cBus *i2c.I2C

func main() {
	machine.Serial.Configure(machine.UARTConfig{
		BaudRate: 115200,
	})

	var rBuf, wBuf [256]byte

	r := bufio.NewReader(&fixReader{Reader: machine.Serial})

	for {
		cmd, err := r.ReadByte()
		if err != nil {
			panic(err)
		}
		if cmd == 0 {
			continue
		}
		switch cmd {
		case 0:
			// do nothing
		case '?':
			io.WriteString(machine.Serial, "OK")
		case 'D': // set pin direction
			applyPins(r, func(p machine.Pin, isInput bool) {
				if isInput {
					p.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
				} else {
					p.Configure(machine.PinConfig{Mode: machine.PinOutput})
				}
			})
		case 'V':
			applyPins(r, func(p machine.Pin, isHigh bool) {
				if isHigh {
					p.High()
				} else {
					p.Low()
				}
			})
		case 'R':
			var v uint16
			for i, p := range pins {
				if p.Get() {
					v |= 1 << uint(i)
				}
			}
			binary.Write(machine.Serial, binary.LittleEndian, v)
		case 'I':
			i2cBus = nil
			var args struct {
				SDA, SCL byte
			}
			if err := binary.Read(r, binary.LittleEndian, &args); err != nil {
				continue
			}
			i2cBus = i2c.New(i2c.NewSoftController(
				driver.FromMachine(pins[args.SDA]),
				driver.FromMachine(pins[args.SCL]),
			))
		case 'i':
			var args struct {
				Addr   uint16
				ReadN  uint8
				WriteN uint8
			}
			if err := binary.Read(r, binary.LittleEndian, &args); err != nil {
				continue
			}
			_, err = io.ReadFull(r, wBuf[:args.WriteN])
			if err != nil {
				continue
			}

			i2cBus.Tx(args.Addr, wBuf[:args.WriteN], rBuf[:args.ReadN])
			machine.Serial.Write(rBuf[:args.ReadN])
		}
	}
}

func applyPins(r io.Reader, f func(machine.Pin, bool)) error {
	var arg uint16
	err := binary.Read(r, binary.LittleEndian, &arg)
	if err != nil {
		return err
	}
	for i, p := range pins {
		f(p, arg&(1<<i) != 0)
	}
	return nil
}
