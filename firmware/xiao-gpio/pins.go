package main

import (
	"errors"
	"machine"
)

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

func setPin(n int, v bool) error {
	if n >= len(pins) {
		return errors.New("xb: invalid pin")
	}
	if n < 0 {
		return errors.New("xb: invalid pin")
	}

	pins[n].Set(v)
	return nil
}

func getPin(n int) (bool, error) {
	if n >= len(pins) {
		return false, errors.New("xb: invalid pin")
	}
	if n < 0 {
		return false, errors.New("xb: invalid pin")
	}

	return pins[n].Get(), nil
}

func setInputPin(n int, v bool) error {
	if n >= len(pins) {
		return errors.New("xb: invalid pin")
	}
	if n < 0 {
		return errors.New("xb: invalid pin")
	}

	if v {
		pins[n].Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	} else {
		pins[n].Configure(machine.PinConfig{Mode: machine.PinOutput})
	}

	return nil
}
