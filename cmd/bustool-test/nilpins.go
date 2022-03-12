package main

type nilPin bool

func (p nilPin) Get() bool { return bool(p) }
func (nilPin) OutputLow()  {}
func (nilPin) PullupHigh() {}
