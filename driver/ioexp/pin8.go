package ioexp

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
	if v {
		*p |= (1 << n)
	} else {
		*p &= ^(1 << n)
	}
}

func (p *Pin8) Toggle(n int) {
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
