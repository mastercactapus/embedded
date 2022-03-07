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

func (Pin16) Len() int { return 8 }

func (p Pin16) Value(n int) bool {
	if n < 0 || n >= 16 {
		return false
	}
	return (p & (1 << n)) != 0
}

func (p *Pin16) High(n ...int) {
	for _, i := range n {
		p.Set(i, true)
	}
}

func (p *Pin16) Low(n ...int) {
	for _, i := range n {
		p.Set(i, false)
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

func (p *Pin16) Toggle(n int) {
	*p ^= (1 << n)
}

func (p *Pin16) ToggleAll() {
	*p = ^*p
}

func (p *Pin16) SetAll(v bool) {
	if v {
		*p = 0xff
	} else {
		*p = 0
	}
}

func (p Pin16) Map(fn func(int) int) PinState {
	if fn == nil {
		return &p
	}

	var n Pin16
	for i := 0; i < 16; i++ {
		n.Set(fn(i), p.Value(i))
	}
	return &n
}
