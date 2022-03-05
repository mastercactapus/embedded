package lcd

import "io"

// The Controller interface is the main interface for interacting with a
// LCD display.
type Controller interface {
	// IsEightBitMode should return true if the controller supports 8-bit communication.
	IsEightBitMode() bool

	// SetBacklight should update the state of the backlight.
	SetBacklight(bool) error

	// WriteByteIR writes to the instruction register.
	WriteByteIR(byte) error

	// ReadByteIR reads from the instruction register.
	ReadByteIR() (byte, error)

	// WriteByte writes to the data register.
	io.ByteWriter

	// ReadByte reads from the data register.
	io.ByteReader
}
