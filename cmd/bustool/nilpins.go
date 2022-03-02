package main

import "github.com/mastercactapus/embedded/i2c"

type nilPin bool

var _ i2c.Pin = nilPin(false)

func (p nilPin) Get() bool { return bool(p) }
func (nilPin) OutputLow()  {}
func (nilPin) PullupHigh() {}
