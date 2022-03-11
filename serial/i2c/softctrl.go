package i2c

import "time"

type Pin interface {
	PullupHigh()
	OutputLow()
	Get() bool
}

type softCtrl struct {
	sda, scl Pin
}

func NewSoftController(scl, sda Pin) Controller {
	return &softCtrl{scl: scl, sda: sda}
}

func (s *softCtrl) wait() {
	time.Sleep(5 * time.Microsecond)
}

func (s *softCtrl) clockUp() {
	s.scl.PullupHigh()
	for !s.scl.Get() {
	}
}

func (s *softCtrl) Start() {
	s.clockUp()
	s.wait()
	s.sda.OutputLow()
	s.wait()
	s.scl.OutputLow()
	s.wait()
}

func (s *softCtrl) Stop() {
	s.scl.PullupHigh()
	s.sda.PullupHigh()
}

func (s *softCtrl) WriteBit(bit bool) {
	if bit {
		s.sda.PullupHigh()
	} else {
		s.sda.OutputLow()
	}
	s.wait()
	s.clockUp()
	s.wait()
	s.scl.OutputLow()
	s.wait()
}

func (s *softCtrl) ReadBit() (value bool) {
	s.sda.PullupHigh()
	s.wait()
	s.clockUp()
	s.wait()
	value = s.sda.Get()
	if !value {
		// keep it low
		s.sda.OutputLow()
	}
	s.wait()
	s.scl.OutputLow()
	s.wait()

	if !value {
		s.sda.PullupHigh()
	}
	return value
}
