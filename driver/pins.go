package driver

import "errors"

type OutputPin interface {
	High() error
	Low() error
	Set(bool) error
}

type InputPin interface {
	Get() (bool, error)
}

type OCPin interface {
	Input() error
	Output() error
	SetInput(bool) error
}

type IOPin interface {
	OutputPin
	InputPin
}

type OCInputPin interface {
	OCPin
	InputPin
}
type Pin interface {
	OCPin
	IOPin
}

type Pinner interface {
	Pin(int) Pin
	PinCount() int
}

type BufferedPinner interface {
	// BufferedPin returns a pin that does not result in a write to the
	// underlying pin until a change from another source is made, or Update is called.
	BufferedPin(int) Pin

	// Flush sends any pending pin state updates to the underlying pin.
	Flush() error

	// Refresh updates the state of all pins, affecting calls to Get().
	Refresh() error

	PinCount() int
}

type PinFN struct {
	N            int
	SetInputFunc func(int, bool) error
	SetFunc      func(int, bool) error
	GetFunc      func(int) (bool, error)
}

func (p PinFN) SetInput(v bool) error {
	if p.SetInputFunc == nil {
		return ErrNotSupported
	}
	return p.SetInputFunc(p.N, v)
}

func (p PinFN) Input() error  { return p.SetInput(true) }
func (p PinFN) Output() error { return p.SetInput(false) }

func (p PinFN) Set(v bool) error {
	if p.SetFunc == nil {
		return ErrNotSupported
	}
	return p.SetFunc(p.N, v)
}
func (p PinFN) High() error { return p.Set(true) }
func (p PinFN) Low() error  { return p.Set(false) }

func (p PinFN) Get() (bool, error) {
	if p.GetFunc == nil {
		return false, ErrNotSupported
	}
	return p.GetFunc(p.N)
}

type PinF struct {
	SetInputFunc func(bool) error
	SetFunc      func(bool) error
	GetFunc      func() (bool, error)
}

var ErrNotSupported = errors.New("not supported")

func (p PinF) SetInput(v bool) error {
	if p.SetInputFunc == nil {
		return ErrNotSupported
	}
	return p.SetInputFunc(v)
}
func (p PinF) Input() error  { return p.SetInput(true) }
func (p PinF) Output() error { return p.SetInput(false) }

func (p PinF) Set(v bool) error {
	if p.SetFunc == nil {
		return ErrNotSupported
	}
	return p.SetFunc(v)
}

func (p PinF) High() error { return p.Set(true) }
func (p PinF) Low() error  { return p.Set(false) }

func (p PinF) Get() (bool, error) {
	if p.GetFunc == nil {
		return false, ErrNotSupported
	}
	return p.GetFunc()
}
