//go:build tinygo
// +build tinygo

package driver

import "machine"

type machinePin machine.Pin

func FromMachine(p machine.Pin) Pin {
	return machinePin(p)
}

func (p machinePin) Get() (bool, error) {
	return machine.Pin(p).Get(), nil
}

func (p machinePin) Set(v bool) error {
	machine.Pin(p).Set(v)
	return nil
}
func (p machinePin) High() error { return p.Set(true) }
func (p machinePin) Low() error  { return p.Set(false) }
func (p machinePin) SetInput(v bool) error {
	if v {
		return p.Input()
	}

	return p.Output()
}

func (p machinePin) Input() error {
	machine.Pin(p).Configure(machine.PinConfig{Mode: machine.PinInput})
	return nil
}

func (p machinePin) Output() error {
	machine.Pin(p).Configure(machine.PinConfig{Mode: machine.PinOutput})
	return nil
}
