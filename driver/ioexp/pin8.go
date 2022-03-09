package ioexp

// PinByte converts a valuer to byte by taking the lowest 8 bits.
func PinByte(v Valuer) (r byte) {
	switch p := v.(type) {
	case Pin8:
		return byte(p)
	case *Pin8:
		return byte(*p)
	}

	for i := 0; i < 8; i++ {
		if v.Value(i) {
			r |= 1 << i
		}
	}
	return r
}

// Pin8 is a PinState for 8-bit expanders.
type Pin8 uint8

func (p Pin8) Value(n int) bool {
	if n < 0 || n >= 8 {
		return false
	}

	return (p & (1 << n)) != 0
}

func (p *Pin8) Set(n int, v bool) {
	if n < 0 || n >= 8 {
		return
	}
	if v {
		*p |= (1 << n)
	} else {
		*p &= ^(1 << n)
	}
}
