//go:build pico
// +build pico

package i2c

import "device"

func wait() {
	for i := 0; i < 50; i++ {
		device.Asm("nop")
	}
}
