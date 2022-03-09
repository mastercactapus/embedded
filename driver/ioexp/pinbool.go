package ioexp

// PinBool stores an arbitrary number of pins as bool.
type PinBool []bool

type pinInvert struct {
	v Valuer
}

func Invert(pins Valuer) Valuer {
	return &pinInvert{pins}
}

func (p *pinInvert) Value(n int) bool { return !p.v.Value(n) }

// NewPins returns a new PinState with the provided number of pins.
func NewPins(n int) PinState {
	switch n {
	case 8:
		return new(Pin8)
	case 16:
		return new(Pin16)
	}

	p := make(PinBool, n)
	return &p
}

func (p PinBool) Len() int { return len(p) }
func (p PinBool) Value(n int) bool {
	if n > len(p) {
		return false
	}
	return p[n]
}

func (p PinBool) Set(n int, v bool) {
	if n > len(p) {
		return
	}
	p[n] = v
}

func (p PinBool) SetAll(v bool) {
	for i := range p {
		p[i] = v
	}
}
