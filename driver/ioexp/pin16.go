package ioexp

// PinByte converts a valuer to byte by taking the lowest 8 bits.
func PinUint16(v Valuer) (r uint16) {
	switch p := v.(type) {
	case Pin16:
		return uint16(p)
	case *Pin16:
		return uint16(*p)
	}

	for i := 0; i < 16; i++ {
		if v.Value(i) {
			r |= 1 << i
		}
	}
	return r
}

// Pin16 is a PinState for 16-bit expanders.
type Pin16 uint16

func (p Pin16) Value(n int) bool {
	if n < 0 || n >= 16 {
		return false
	}
	return (p & (1 << n)) != 0
}

func SetHigh(s PinState, n ...int) {
	for _, i := range n {
		s.Set(i, true)
	}
}

func SetLow(s PinState, n ...int) {
	for _, i := range n {
		s.Set(i, false)
	}
}

func Toggle(s PinState, n ...int) {
	for _, i := range n {
		s.Set(i, !s.Value(i))
	}
}

func (p *Pin16) Set(n int, v bool) {
	if n < 0 || n >= 16 {
		return
	}
	if v {
		*p |= (1 << n)
	} else {
		*p &= ^(1 << n)
	}
}
