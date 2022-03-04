package ioexp

// Pin8 is a PinState for 8-bit expanders.
type Pin8 uint8

func (Pin8) Len() int { return 8 }

func (p Pin8) Get(n int) bool {
	if n < 0 || n >= 8 {
		return false
	}

	return (p & (1 << n)) != 0
}

// Sets the numbered pin to HIGH/Input for true or LOW/Output for false.
func (p *Pin8) Set(n int, v bool) {
	if v {
		*p |= (1 << n)
	} else {
		*p &= ^(1 << n)
	}
}
