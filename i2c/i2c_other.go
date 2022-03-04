//go:build !pico
// +build !pico

package i2c

import "time"

func wait() {
	time.Sleep(10 * time.Microsecond)
}
