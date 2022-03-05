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

func (Pin8) Len() int { return 8 }

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

func (p *Pin8) Toggle(n int) {
	if n < 0 || n >= 8 {
		return
	}
	*p ^= (1 << n)
}

func (p *Pin8) ToggleAll() {
	*p = ^*p
}

func (p *Pin8) SetAll(v bool) {
	if v {
		*p = 0xff
	} else {
		*p = 0
	}
}

func (p Pin8) Map(fn func(int) int) PinState {
	if fn == nil {
		return &p
	}

	var n Pin8
	for i := 0; i < 8; i++ {
		n.Set(fn(i), p.Value(i))
	}
	return &n
}
