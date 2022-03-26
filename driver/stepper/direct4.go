package stepper

import "github.com/mastercactapus/embedded/driver"

func New4Phase(a, b, c, d driver.OutputPin) *Direct {
	return &Direct{
		Sequence: []byte{
			0b1000,
			0b0100,
			0b0010,
			0b0001,
		},
		Pins: []driver.OutputPin{a, b, c, d},
		N:    -1,
	}
}

type Direct struct {
	Sequence []byte
	Pins     []driver.OutputPin
	Reverse  bool

	N int
}

func (d *Direct) Off() error {
	for _, p := range d.Pins {
		err := p.Set(false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Direct) On() (err error) {
	for i, p := range d.Pins {
		err = p.Set(d.Sequence[i]>>uint(d.N)&1 != 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Direct) Step() error {
	if d.Reverse {
		if d.N == 0 {
			d.N = len(d.Sequence) - 1
		} else {
			d.N--
		}
	} else {
		d.N++
		if d.N == len(d.Sequence) {
			d.N = 0
		}
	}
	return d.On()
}
