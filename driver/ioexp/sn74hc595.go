package ioexp

import (
	"github.com/mastercactapus/embedded/driver"
)

type SN74HC595 struct {
	cfg SN74HC595Config

	state []uint8
	clean bool
	err   error
}

type SN74HC595Config struct {
	SRCLK driver.OutputPin
	RCLK  driver.OutputPin
	SER   driver.OutputPin

	// Optionally set the clear pin before writing data.
	SRCLR driver.OutputPin
}

func NewSN74HC595(cfg SN74HC595Config) *SN74HC595 {
	return &SN74HC595{
		cfg: cfg,
	}
}

// Configure sets the starting state of the shift register.
//
// If there are multiple registers chained together, the first
// state value should be the first register in the chain.
func (s *SN74HC595) Configure(state ...uint8) error {
	s.state = make([]uint8, len(state))
	// copy in reverse order
	for i, v := range state {
		s.state[len(state)-1-i] = v
	}

	return s.Flush()
}

func (s *SN74HC595) PinCount() int { return len(s.state) * 8 }

func (s *SN74HC595) Pin(n int) driver.Pin {
	return &driver.PinFN{
		N:       n,
		SetFunc: s.setPin,
	}
}

func (s *SN74HC595) BufferedPin(n int) driver.Pin {
	return &driver.PinFN{
		N:       n,
		SetFunc: s.setPinBuf,
	}
}

func (s *SN74HC595) setPinBuf(n int, v bool) error {
	if v {
		s.state[n/8] |= 1 << uint(n%8)
	} else {
		s.state[n/8] &= ^(1 << uint(n%8))
	}
	return nil
}

func (s *SN74HC595) setPin(n int, v bool) error {
	s.setPinBuf(n, v)
	return s.Flush()
}

func (s *SN74HC595) Flush() error {
	if s.cfg.SRCLR != nil {
		s.pulse(s.cfg.SRCLR)
	} else if !s.clean {
		s.setLow(s.cfg.SER)
		for i := s.PinCount(); i > 0; i-- {
			s.pulse(s.cfg.SRCLK)
		}
		s.clean = true
	}

	for _, b := range s.state {
		for i := uint(0); i < 8; i++ {
			s.writeBit(b&(1<<i) != 0)
		}
	}

	s.pulse(s.cfg.RCLK)

	err := s.err
	s.err = nil
	s.clean = err == nil
	return err
}

func (s *SN74HC595) setHigh(p driver.OutputPin) {
	if s.err != nil {
		return
	}
	s.err = p.High()
}

func (s *SN74HC595) setLow(p driver.OutputPin) {
	if s.err != nil {
		return
	}
	s.err = p.Low()
}

func (s *SN74HC595) pulse(p driver.OutputPin) {
	if s.err != nil {
		return
	}
	s.setHigh(p)
	s.setLow(p)
}

func (s *SN74HC595) writeBit(val bool) {
	if val {
		s.setHigh(s.cfg.SER)
	} else {
		s.setLow(s.cfg.SER)
	}
	s.pulse(s.cfg.SRCLK)
}
