package spi

import (
	"github.com/mastercactapus/embedded/driver"
)

type Config struct {
	Mode

	SCLK, MOSI, MISO driver.Pin

	Baud int
}

type SoftCtrl struct {
	*Config
	err  error
	fill byte
	wait func()
}

var (
	_ Controller          = (*SoftCtrl)(nil)
	_ ReadWriteController = (*SoftCtrl)(nil)
	_ ReadController      = (*SoftCtrl)(nil)
	_ WriteController     = (*SoftCtrl)(nil)
)

func NewSoftCtrl(cfg *Config) (*SoftCtrl, error) {
	if cfg.Baud == 0 {
		cfg.Baud = 100000
	}
	s := &SoftCtrl{
		Config: cfg,
		wait:   func() {},
	}

	s.clockIdle()
	if err := s.MOSI.Output(); err != nil {
		return nil, err
	}
	if err := s.SCLK.Output(); err != nil {
		return nil, err
	}
	if err := s.MISO.Input(); err != nil {
		return nil, err
	}
	if err := s.MOSI.High(); err != nil {
		return nil, err
	}
	s.wait()
	return s, s.readErr()
}

func (s *SoftCtrl) readErr() (err error) {
	err = s.err
	s.err = nil
	return err
}

func (s *SoftCtrl) clockIdle() {
	if s.err != nil {
		return
	}

	if s.CPOL() {
		s.err = s.SCLK.High()
	} else {
		s.err = s.SCLK.Low()
	}
}

func (s *SoftCtrl) clockActive() {
	if s.err != nil {
		return
	}

	if s.CPOL() {
		s.err = s.SCLK.Low()
	} else {
		s.err = s.SCLK.High()
	}
}

func (s *SoftCtrl) setMOSI(v bool) {
	if s.err != nil {
		return
	}

	s.err = s.MOSI.Set(v)
}

func (s *SoftCtrl) readMISO() (v bool) {
	if s.err != nil {
		return false
	}

	v, s.err = s.MISO.Get()
	return v
}

func (s *SoftCtrl) Write(p []byte) (int, error) {
	for i := range p {
		_, s.err = s.ReadWriteByte(p[i])
		if s.err != nil {
			return i, s.readErr()
		}
	}
	return len(p), s.readErr()
}
func (s *SoftCtrl) SetFill(fill byte) error { s.fill = fill; return nil }
func (s *SoftCtrl) ReadByte() (byte, error) { return s.ReadWriteByte(s.fill) }
func (s *SoftCtrl) Read(p []byte) (int, error) {
	for i := range p {
		p[i], s.err = s.ReadWriteByte(s.fill)
		if s.err != nil {
			return i, s.readErr()
		}
	}
	return len(p), s.readErr()
}

func (s *SoftCtrl) ReadWrite(p []byte) (int, error) {
	for i := range p {
		p[i], s.err = s.ReadWriteByte(p[i])
		if s.err != nil {
			return i, s.readErr()
		}
	}
	return len(p), s.readErr()
}

func (s *SoftCtrl) ReadWriteByte(v byte) (byte, error) {
	if s.CPHA() {
		for i := 7; i >= 0; i-- {
			s.clockActive()
			s.setMOSI((v & (1 << uint(i))) != 0)
			s.wait()
			s.clockIdle()
			if s.readMISO() {
				v |= 1 << uint(i)
			} else {
				v &= ^(1 << uint(i))
			}
			s.wait()
		}
	} else {
		for i := 0; i < 8; i++ {
			s.setMOSI((v & 0x80) != 0)
			v <<= 1
			s.clockActive()
			if s.readMISO() {
				v |= 1
			} else {
				v &= 0xfe
			}
			s.wait()
			s.clockIdle()
			s.wait()
		}
	}

	return v, s.readErr()
}
