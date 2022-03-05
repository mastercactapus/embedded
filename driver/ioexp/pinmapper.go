package ioexp

type PinMapper struct {
	w     PinWriter
	r     PinReader
	wMap  func(int) int
	rMap  func(int) int
	count int
}

func (p *PinMapper) PinCount() int { return p.count }

// NewPinMapper returns a PinMapper that maps all pins being read and written
// with the given functions.
//
// writeMap and readMap will map pins passed to WritePins and ReadPins, respectively.
func NewPinMapper(rw PinReadWriter, writeMapFn, readMapFn func(int) int) PinReadWriter {
	return &PinMapper{
		w:     rw,
		r:     rw,
		wMap:  writeMapFn,
		rMap:  readMapFn,
		count: rw.PinCount(),
	}
}

// NewPinWriterMapper returns a PinWriter that maps all pins being set with
// the given function.
func NewPinWriterMapper(w PinWriter, mapFn func(int) int) PinWriter {
	return &PinMapper{
		w:     w,
		wMap:  mapFn,
		count: w.PinCount(),
	}
}

func (p *PinMapper) WritePins(v Valuer) error {
	return p.w.WritePins(v.Map(p.wMap))
}

// NewPinReaderMapper returns a PinReader that maps all pins being read with
// the given function.
func NewPinReaderMapper(r PinReader, mapFn func(int) int) PinReader {
	return &PinMapper{
		r:     r,
		rMap:  mapFn,
		count: r.PinCount(),
	}
}

func (p *PinMapper) ReadPins() (PinState, error) {
	v, err := p.r.ReadPins()
	if err != nil {
		return nil, err
	}
	return v.Map(p.rMap), nil
}
