//go:build pico
// +build pico

package main

import "machine"

func configIO() (io.Reader, io.Writer) {
	err := machine.Serial.Configure(machine.UARTConfig{
		BaudRate: 115200,
	})
	if err != nil {
		panic(err)
	}
	return machine.Serial, machine.Serial
}
