package ioexp

// An PinReadWriter is a generic IO expander interface.
type PinReadWriter interface {
	// WritePins will update all IO pins to the given state.
	WritePins(PinState) error

	// ReadPins will read the current *VALUE* of all IO pins.
	//
	// This may be different than the last written state.
	ReadPins() (PinState, error)

	// PinCount returns the number of pins available on the device.
	PinCount() int
}
