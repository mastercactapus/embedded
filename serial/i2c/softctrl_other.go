//go:build !tinygo
// +build !tinygo

package i2c

import "time"

func (s *softCtrl) wait() { time.Sleep(5 * time.Microsecond) }
