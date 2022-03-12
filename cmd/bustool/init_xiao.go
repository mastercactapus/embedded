//go:build xiao
// +build xiao

package main

import (
	"io"
	"machine"
)

func configIO() (io.Reader, io.Writer) {
	machine.Serial.Configure(machine.UARTConfig{
		BaudRate: 115200,
	})

	return machine.Serial, machine.Serial
}
