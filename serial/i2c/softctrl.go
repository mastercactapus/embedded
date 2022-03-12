package i2c

import (
	"github.com/mastercactapus/embedded/driver"
)

type softCtrl struct {
	sda, scl driver.Pin
	err      error
}

// NewSoftController will create a generic I2C controller.
// Pins should be configured as inputs with pull-up and output
// should be logic-low.
func NewSoftController(sda, scl driver.Pin) Controller {
	sda.High()
	sda.Input()
	scl.High()
	scl.Input()
	return &softCtrl{scl: scl, sda: sda}
}

func (s *softCtrl) setHigh(p driver.Pin) {
	if s.err != nil {
		return
	}

	s.err = p.Input()
	if s.err != nil {
		return
	}
	s.err = p.High()
}

func (s *softCtrl) setLow(p driver.Pin) {
	if s.err != nil {
		return
	}
	s.err = p.Output()
	if s.err != nil {
		return
	}
	s.err = p.Low()
}

func (s *softCtrl) get(p driver.Pin) (val bool) {
	if s.err != nil {
		return false
	}

	val, s.err = p.Get()
	return val
}

func (s *softCtrl) waitHigh(p driver.Pin) {
	if s.err != nil {
		return
	}
	s.setHigh(p)
	if s.err != nil {
		return
	}
	var v bool
	// TODO: timeout
	for {
		v, s.err = p.Get()
		if s.err != nil {
			return
		}
		if v {
			return
		}
	}
}

func (s *softCtrl) clockUp() { s.waitHigh(s.scl) }

func (s *softCtrl) Start() {
	s.clockUp()
	s.wait()
	s.setLow(s.sda)
	s.wait()
	s.setLow(s.scl)
	s.wait()
}

func (s *softCtrl) Stop() {
	s.clockUp()
	s.waitHigh(s.sda)
}

func (s *softCtrl) WriteBit(bit bool) {
	if bit {
		s.setHigh(s.sda)
	} else {
		s.setLow(s.sda)
	}
	s.wait()
	s.clockUp()
	s.wait()
	s.setLow(s.scl)
	s.wait()
}

func (s *softCtrl) ReadBit() (value bool) {
	s.setHigh(s.sda)
	s.wait()
	s.clockUp()
	s.wait()
	value = s.get(s.sda)
	if !value {
		// keep it low
		s.setLow(s.sda)
	}
	s.wait()
	s.setLow(s.scl)
	s.wait()

	if !value {
		s.setHigh(s.sda)
	}
	return value
}
