package ioexp

type PinWriter interface {
	// WritePins will update all IO pins to the given state.
	//
	// If the device does not support writing, ErrReadOnly should be returned.
	//
	// Note: Some devices will swap input/output modes based on the pin state.
	WritePins(Valuer) error

	// WritePinsMask works like WritePins but only updates the masked pins.
	//
	// Note: Some devices will update all pins on first WritePinsMask call.
	WritePinsMask(pins, mask Valuer) error

	// PinCount returns the number of pins available on the device.
	PinCount() int
}

type InputSetter interface {
	// SetInputPins should set all pins returning true to input,
	// and the rest as output.
	SetInputPins(Valuer) error

	// SetInputPinsMask works like SetInputPins but only updates the masked pins.
	SetInputPinsMask(pins, mask Valuer) error
}
type PullupSetter interface {
	// SetPullupPins should enable pullups on all pins returning true.
	//
	// The caller should be careful to only enable pullups on pins that
	// are actually inputs. Some implementations may change the pin mode
	// to input if the pullup is enabled.
	SetPullupPins(Valuer) error

	// SetPullupPinsMask works like SetPullupPins but only updates the masked pins.
	SetPullupPinsMask(pins, mask Valuer) error
}

type InvertSetter interface {
	// SetInvertPins will cause the pin value from ReadPins to be inverted
	// for all pins returning true.
	SetInvertPins(Valuer) error

	// SetInvertPinsMask works like SetInvertPins but only updates the masked pins.
	SetInvertPinsMask(pins, mask Valuer) error
}

type PinReader interface {
	// ReadPins will read the current *VALUE* of all IO pins.
	//
	// Implementers should ensure the value is that of the actual pin state, and not
	// the register state.
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
