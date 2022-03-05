package ioexp

type PinWriter interface {
	// WritePins will update all IO pins to the given state.
	//
	// If the device does not support writing, ErrReadOnly should be returned.
	WritePins(Valuer) error

	// PinCount returns the number of pins available on the device.
	PinCount() int
}

type PinReader interface {
	// ReadPins will read the current *VALUE* of all IO pins.
	//
	// This may be different than the last written state.
	//
	// If the device does not support reading, ErrWriteOnly should be returned.
	ReadPins() (PinState, error)

	// PinCount returns the number of pins available on the device.
	PinCount() int
}

// An PinReadWriter is a generic IO expander interface.
type PinReadWriter interface {
	PinWriter
	PinReader
}
