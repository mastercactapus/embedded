package ioexp

type Valuer interface {
	// Value returns true if the pin is HIGH, false if LOW.
	//
	// Calling Get on a pin that is not available will always return false.
	Value(int) bool

	// Len returns the number of underlying pins.
	Len() int

	// Map returns a new PinState that is the result of applying the given
	// function to each pin.
	//
	// The returned PinState will be the same type as the receiver and
	// is guaranteed to have the same number of pins.
	//
	// Map may be called with nil to return a copy of the pins.
	Map(func(int) int) PinState
}

// PinState is a generic interface for IO expander pin state.
//
// Updating values should not control actual pins, it should
// be passed to the PinWriter interface.
type PinState interface {
	Valuer
	Setter
}

type Setter interface {
	// Set sets the numbered pin to HIGH for true or LOW for false.
	//
	// Calling Set on a pin that is not available will be ignored.
	Set(int, bool)

	// Toggle swaps the value of the numbered pin.
	//
	// Calling Toggle on a pin that is not available will be ignored.
	Toggle(int)

	// SetAll is a convenience method that sets all pins to the provided value.
	SetAll(bool)

	// ToggleAll is a convenience method that toggles all pins.
	ToggleAll()

	// Len returns the number of underlying pins.
	Len() int
}

// CopyState copies the state of one PinState to another.
func CopyState(dst Setter, src Valuer) {
	if dt, ok := dst.(*Pin8); ok {
		switch s := src.(type) {
		case *Pin8:
			*dt = *s
			return
		case Pin8:
			*dt = s
			return
		}
	}

	n := dst.Len()
	if n > src.Len() {
		n = src.Len()
	}
	for i := 0; i < n; i++ {
		dst.Set(i, src.Value(i))
	}
}
