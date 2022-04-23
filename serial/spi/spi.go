package spi

type Controller interface {
	ReadWriteByte(byte) (byte, error)
}

type ReadWriteController interface {
	ReadWrite([]byte) (int, error)
}

type ReadController interface {
	SetFill(byte) error
	Read([]byte) (int, error)
}

type WriteController interface {
	Write([]byte) (int, error)
}

func ReadWrite(c Controller, p []byte) (n int, err error) {
	if rw, ok := c.(ReadWriteController); ok {
		return rw.ReadWrite(p)
	}

	for i := range p {
		p[i], err = c.ReadWriteByte(p[i])
		if err != nil {
			return i, err
		}
	}
	return len(p), nil
}

func Read(c Controller, fill byte, p []byte) (n int, err error) {
	if r, ok := c.(ReadController); ok {
		if err := r.SetFill(fill); err != nil {
			return 0, err
		}
		return r.Read(p)
	}

	if rw, ok := c.(ReadWriteController); ok {
		for i := range p {
			p[i] = fill
		}
		return rw.ReadWrite(p)
	}

	for i := range p {
		p[i], err = c.ReadWriteByte(fill)
		if err != nil {
			return i, err
		}
	}
	return len(p), nil
}

func Write(c Controller, p []byte) (n int, err error) {
	if w, ok := c.(WriteController); ok {
		return w.Write(p)
	}

	if rw, ok := c.(ReadWriteController); ok {
		buf := make([]byte, len(p))
		copy(buf, p)
		return rw.ReadWrite(buf)
	}

	for i := range p {
		_, err = c.ReadWriteByte(p[i])
		if err != nil {
			return i, err
		}
	}

	return len(p), nil
}

type SPI struct {
	Controller
	f byte
}

func (s *SPI) SetFillByte(f byte) { s.f = f }

func (s *SPI) Read(p []byte) (int, error) {
	var err error
	for i := range p {
		p[i], err = s.ReadWriteByte(s.f)
		if err != nil {
			return i, err
		}
	}
	return len(p), nil
}

func (s *SPI) Write(p []byte) (int, error) {
	var err error
	for i := range p {
		_, err = s.ReadWriteByte(p[i])
		if err != nil {
			return i, err
		}
	}
	return len(p), nil
}

// func (s *SPI) Tx(w, r []byte) (err error) {
// 	var i int
// 	for {
// 		switch {
// 		case i < len(w) && i < len(r):
// 			r[i], err = s.ReadWriteByte(w[i])
// 		case i < len(w):
// 			_, err = s.ReadWriteByte(w[i])
// 		case i < len(r):
// 			r[i], err = s.ReadWriteByte(0)
// 		default:
// 			return nil
// 		}
// 		if err != nil {
// 			return err
// 		}
// 	}
// }
