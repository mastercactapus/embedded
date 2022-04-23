package spi

type Mode byte

const (
	Mode0 = Mode(0b00)
	Mode1 = Mode(0b01)
	Mode2 = Mode(0b10)
	Mode3 = Mode(0b11)
)

// CPOL returns the clock polarity.
//
// If true, the clock is idle HIGH and each
// cycle the clock is toggled LOW.
func (m Mode) CPOL() bool { return m&0b10 != 0 }

// CPHA returns the clock phase.
//
// If true, data is sampled on the leading edge
// of the clock and changed on the trailing edge.
//
// If false, data is sampled on the trailing edge
// and changed on the leading edge.
func (m Mode) CPHA() bool { return m&0b01 != 0 }
