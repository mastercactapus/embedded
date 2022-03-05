package ioexp

type funcWriter struct {
	fn    func(int, bool) error
	count int
}

// NewFuncWriter returns a PinWriter that calls fn(pin, val) on every WritePins call.
func NewFuncWriter(fn func(int, bool) error, n int) PinWriter {
	return &funcWriter{fn: fn}
}

func (fw *funcWriter) WritePins(pin Valuer) error {
	for i := 0; i < pin.Len(); i++ {
		if err := fw.fn(i, pin.Value(i)); err != nil {
			return err
		}
	}
	return nil
}
func (fw *funcWriter) PinCount() int { return fw.count }
