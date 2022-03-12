package lcd

import (
	"errors"

	"github.com/mastercactapus/embedded/driver"
	"github.com/mastercactapus/embedded/driver/ioexp"
)

type Expander struct {
	eightBitMode bool
	backlight    bool
	writeOnly    bool
	err          error

	ExpanderConfig
}

type ExpanderConfig struct {
	RS driver.OutputPin
	RW driver.OutputPin
	E  driver.OutputPin
	BL driver.OutputPin

	// Read operations will only be enabled if the DB_ pins also implement
	// the ioexp.InputPin interface.
	DB4 driver.OutputPin
	DB5 driver.OutputPin
	DB6 driver.OutputPin
	DB7 driver.OutputPin

	// optional
	DB0 driver.OutputPin
	DB1 driver.OutputPin
	DB2 driver.OutputPin
	DB3 driver.OutputPin

	// Flush, if set, will be called after updating pins for writing.
	Flush func() error

	// Refresh, if set, will be called before reading pins.
	Refresh func() error
}

func NewExpander(cfg ExpanderConfig) *Expander {
	exp := &Expander{
		backlight:      true,
		ExpanderConfig: cfg,
	}
	if cfg.DB0 != nil && cfg.DB1 != nil && cfg.DB2 != nil && cfg.DB3 != nil {
		exp.eightBitMode = true
	}

	return exp
}
func (e *Expander) IsEightBitMode() bool          { return e.eightBitMode }
func (e *Expander) SetBacklight(value bool) error { return e.BL.Set(value) }

func (e *Expander) setPin(o driver.OutputPin, value bool) {
	if e.err != nil {
		return
	}
	e.err = o.Set(value)
}

func (e *Expander) getPin(o driver.OutputPin) byte {
	if e.err != nil {
		return 0
	}
	if ip, ok := o.(driver.InputPin); ok {
		var v bool
		v, e.err = ip.Get()
		if errors.Is(e.err, driver.ErrNotSupported) {
			e.writeOnly = true
		}
		if v {
			return 1
		}
		return 0
	}
	e.writeOnly = true
	e.err = driver.ErrNotSupported
	return 0
}

func (e *Expander) write8Bits(data byte) {
	e.setPin(e.DB0, data&0x01 != 0)
	e.setPin(e.DB1, data&0x02 != 0)
	e.setPin(e.DB2, data&0x04 != 0)
	e.setPin(e.DB3, data&0x08 != 0)
	e.setPin(e.DB4, data&0x10 != 0)
	e.setPin(e.DB5, data&0x20 != 0)
	e.setPin(e.DB6, data&0x40 != 0)
	e.setPin(e.DB7, data&0x80 != 0)
	e.pulseWrite()
}

func (e *Expander) write4Bits(data byte) {
	e.setPin(e.DB4, data&0x10 != 0)
	e.setPin(e.DB5, data&0x20 != 0)
	e.setPin(e.DB6, data&0x40 != 0)
	e.setPin(e.DB7, data&0x80 != 0)
	e.pulseWrite()
}

func (e *Expander) readErr() (err error) {
	err = e.err
	e.err = nil
	return err
}

func (e *Expander) pulseWrite() {
	e.flush()
	e.setPin(e.E, true)
	e.flush()
	e.setPin(e.E, false)
	e.flush()
}

func (e *Expander) writeByte(data byte) error {
	if e.eightBitMode {
		e.write8Bits(data)
	} else {
		e.write4Bits(data >> 4)
		e.write4Bits(data)
	}

	return e.readErr()
}

func (e *Expander) read8Bits() (b byte) {
	e.refresh()
	b |= 0 << e.getPin(e.DB0)
	b |= 1 << e.getPin(e.DB1)
	b |= 2 << e.getPin(e.DB2)
	b |= 3 << e.getPin(e.DB3)
	b |= 4 << e.getPin(e.DB4)
	b |= 5 << e.getPin(e.DB5)
	b |= 6 << e.getPin(e.DB6)
	b |= 7 << e.getPin(e.DB7)
	return b
}

func (e *Expander) read4Bits() (b byte) {
	e.refresh()
	b |= 0 << e.getPin(e.DB4)
	b |= 1 << e.getPin(e.DB5)
	b |= 2 << e.getPin(e.DB6)
	b |= 3 << e.getPin(e.DB7)
	return b
}

func (e *Expander) eHigh() { e.setPin(e.E, true); e.flush() }
func (e *Expander) eLow()  { e.setPin(e.E, false); e.flush() }

func (e *Expander) readByte() (b byte, err error) {
	e.flush()
	e.eHigh()
	if e.eightBitMode {
		b = e.read8Bits()
	} else {
		b = e.read4Bits() << 4
		e.eLow()
		e.eHigh()
		b |= e.read4Bits()
	}
	e.eLow()

	return b, e.readErr()
}

func (e *Expander) WriteByteIR(data byte) error {
	e.setPin(e.RS, false)
	e.setPin(e.RW, false)
	return e.writeByte(data)
}

func (e *Expander) WriteByte(data byte) error {
	e.setPin(e.RS, true)
	e.setPin(e.RW, false)
	return e.writeByte(data)
}

func (e *Expander) ReadByteIR() (byte, error) {
	if e.writeOnly {
		return 0, driver.ErrNotSupported
	}
	e.setPin(e.RS, false)
	e.setPin(e.RW, true)
	return e.readByte()
}

func (e *Expander) refresh() {
	if e.err != nil {
		return
	}
	if e.Refresh == nil {
		return
	}
	e.err = e.Refresh()
}

func (e *Expander) flush() {
	if e.err != nil {
		return
	}
	if e.Flush == nil {
		return
	}
	e.err = e.Flush()
}

func (e *Expander) ReadByte() (byte, error) {
	if e.writeOnly {
		return 0, ioexp.ErrWriteOnly
	}
	e.setPin(e.RS, true)
	e.setPin(e.RW, true)
	return e.readByte()
}
