package main

import (
	"machine"

	"github.com/mastercactapus/embedded/xb"
)

func main() {
	machine.Serial.Configure(machine.UARTConfig{
		BaudRate: 115200,
	})

	panic(xb.NewServer(&serialReader{}, machine.Serial, &xiao{}).Serve())
}
