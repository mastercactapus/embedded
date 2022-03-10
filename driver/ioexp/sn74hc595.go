package ioexp

type SN74HC595 struct {
	w    PinWriter
	pins PinBool

	m *PinMasker
}

var _ PinWriter = (*SN74HC595)(nil)

// NewSN74HC595 creates a PinWriter that writes to a 74HC595 shift register.
//
// The provided PinWriter should control the following pins:
// 0: SRCLK (shift register clock)
// 1: RCLK (storage register clock/latch)
// 2: SER (serial data)
func NewSN74HC595(w PinWriter, bits int) *SN74HC595 {
	return &SN74HC595{
		w:    w,
		pins: make(PinBool, 3),
		m:    NewPinMasker(bits),
	}
}

func (s *SN74HC595) PinCount() int { return s.pins.Len() }

func (s *SN74HC595) writeBit(val bool) (err error) {
	s.pins.Set(2, val)
	if err = s.w.WritePins(s.pins); err != nil {
		return err
	}
	s.pins.Set(0, true)
	if err = s.w.WritePins(s.pins); err != nil {
		return err
	}
	s.pins.Set(0, false)
	if err = s.w.WritePins(s.pins); err != nil {
		return err
	}
	return nil
}

func (s *SN74HC595) WritePins(pins Valuer) (err error) {
	s.pins.SetAll(false)
	for i := s.PinCount() - 1; i >= 0; i-- {
		if err = s.writeBit(pins.Value(i)); err != nil {
			return err
		}
	}
	s.pins.SetAll(false)
	s.pins.Set(1, true)
	if err = s.w.WritePins(s.pins); err != nil {
		return err
	}
	s.pins.Set(1, false)
	return s.w.WritePins(s.pins)
}

func (s *SN74HC595) WritePinsMask(pins, mask Valuer) error {
	return s.m.ApplyFn(pins, mask, s.WritePins)
}
