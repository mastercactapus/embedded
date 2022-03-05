package ioexp

import "fmt"

type PinMask []int

func (PinMask) Len() int { return -1 }
func (m PinMask) Value(i int) bool {
	for _, v := range m {
		if v == i {
			return true
		}
	}
	return false
}

type PinMasker struct {
	v PinState
	n int
}

func ClonePins(pins Valuer, n int) PinState {
	v := make(PinBool, n)
	for i := 0; i < n; i++ {
		v.Set(i, pins.Value(i))
	}
	return v
}

// NewPinMasker returns a PinMasker with the provided intial value and pin count.
func NewPinMasker(pinCount int) *PinMasker {
	return &PinMasker{n: pinCount}
}

// Set will update the state of all pins.
//
// This can be used from a WritePins method.
func (p *PinMasker) Set(v Valuer) { p.v = ClonePins(v, p.n) }

// Apply returns a copy of the pin mask with the given value applied.
func (p *PinMasker) Apply(pins, mask Valuer) Valuer { return ApplyPinsMaskN(p.v, pins, mask, p.n) }

func (p *PinMasker) ApplyFn(pins, mask Valuer, fn func(Valuer) error) error {
	if p.v == nil {
		return fmt.Errorf("mask: tried to apply mask before writing initial state")
	}

	val := p.Apply(pins, mask)
	err := fn(val)
	if err != nil {
		return err
	}

	p.Set(val)
	return nil
}

// ApplyPinsMaskN applies the given value to the given `n` pins with the provided mask.
func ApplyPinsMaskN(oldPins, newPins, mask Valuer, n int) Valuer {
	v := ClonePins(oldPins, n)

	for i := 0; i < n; i++ {
		if !mask.Value(i) {
			continue
		}
		v.Set(i, newPins.Value(i))
	}

	return v
}
