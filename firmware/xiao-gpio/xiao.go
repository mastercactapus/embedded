package main

import (
	"github.com/mastercactapus/embedded/driver"
)

type xiao struct{}

func (xiao) Pin(n int) driver.Pin {
	return &driver.PinFN{
		N:            n,
		SetInputFunc: setInputPin,
		SetFunc:      setPin,
		GetFunc:      getPin,
	}
}
func (xiao) PinCount() int { return len(pins) }
