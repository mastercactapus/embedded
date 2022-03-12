//go:build tinygo
// +build tinygo

package i2c

import "device"

func (s *softCtrl) wait() {
	wait := 20
	for i := 0; i < wait; i++ {
		device.Asm(`nop`)
	}
}
