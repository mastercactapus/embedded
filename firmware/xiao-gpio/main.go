package main

import (
	"bufio"
	"io"
	"machine"
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

func main() {
	machine.Serial.Configure(machine.UARTConfig{
		BaudRate: 115200,
	})

	machine.Serial.ReadByte()
	r := bufio.NewReader(&fixReader{Reader: machine.Serial})
	for {
		cmd, err := r.ReadByte()
		if err != nil {
			panic(err)
		}

		switch cmd {
		case '?':
			io.WriteString(machine.Serial, "OK")
		case 'D': // set direction
			applyPins(r, func(p machine.Pin, isInput bool) {
				if isInput {
					p.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
				} else {
					p.Configure(machine.PinConfig{Mode: machine.PinOutput})
				}
			})
		case 'V': // set value
			applyPins(r, func(p machine.Pin, isHigh bool) {
				if isHigh {
					p.High()
				} else {
					p.Low()
				}
			})
		case 'R': // read value
			var v uint16
			for i, p := range pins {
				if p.Get() {
					v |= 1 << uint(i)
				}
			}
			machine.Serial.WriteByte(byte(v))
			machine.Serial.WriteByte(byte(v >> 8))
		}
	}
}

func applyPins(r *bufio.Reader, f func(machine.Pin, bool)) {
	var buf [2]byte
	io.ReadFull(r, buf[:])

	var arg uint16
	arg = uint16(buf[0]) | uint16(buf[1])<<8
	for i, p := range pins {
		f(p, arg&(1<<i) != 0)
	}
}
