package ioexp

// PinBool stores an arbitrary number of pins as bool.
type PinBool []bool

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

func (p PinBool) Toggle(n int) {
	if n > len(p) {
		return
	}
	p[n] = !p[n]
}

func (p PinBool) ToggleAll() {
	for i := range p {
		p[i] = !p[i]
	}
}

func (p PinBool) SetAll(v bool) {
	for i := range p {
		p[i] = v
	}
}

func (p PinBool) Map(fn func(int) int) PinState {
	if fn == nil {
		return &p
	}

	n := make(PinBool, len(p))
	for i, v := range p {
		n[fn(i)] = v
	}

	return &n
}
