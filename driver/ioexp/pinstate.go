package ioexp

// PinState is a generic interface for IO expander pin state.
type PinState interface {
	// Get returns true if the pin is HIGH, false if LOW.
	//
	// Calling Get on a pin that is not available will always return false.
	Get(int) bool

	// Set sets the numbered pin to HIGH for true or LOW for false.
	//
	// Calling Set on a pin that is not available will be ignored.
	Set(int, bool)

	// Len returns the number of pins PinState represents.
	Len() int
}

// CopyState copies the state of one PinState to another.
func CopyState(dst, src PinState) {
	dt, dok := dst.(*Pin8)
	st, sok := src.(*Pin8)
	if dok && sok {
		*dt = *st
		return
	}

	n := dst.Len()
	if n > src.Len() {
		n = src.Len()
	}
	for i := 0; i < n; i++ {
		dst.Set(i, src.Get(i))
	}
}
