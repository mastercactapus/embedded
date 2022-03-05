package lcd

import (
	"errors"

	"github.com/mastercactapus/embedded/driver/ioexp"
)

const (
	RS = 0
	RW = 1
	E  = 2

	BL = 3

	DB4 = 4
	DB5 = 5
	DB6 = 6
	DB7 = 7

	DB0 = 8
	DB1 = 7
	DB2 = 10
	DB3 = 11
)

type Expander struct {
	w ioexp.PinWriter

	eightBitMode bool
	backlight    bool
	writeOnly    bool

	pins ioexp.PinState
}

func NewExpander(w ioexp.PinWriter) *Expander {
	if w.PinCount() < 8 {
		panic("lcd: ioexpander must have at least 8 pins")
	}
	return &Expander{
		w: w,

		eightBitMode: w.PinCount() >= 12,
		backlight:    true,

		pins: make(ioexp.PinBool, w.PinCount()),
	}
}
func (e *Expander) IsEightBitMode() bool { return e.eightBitMode }
func (e *Expander) SetBacklight(value bool) error {
	e.pins.Set(BL, value)
	return e.w.WritePins(e.pins)
}

func (e *Expander) write8Bits(data byte) error {
	e.pins.Set(DB0, data&0x01 != 0)
	e.pins.Set(DB1, data&0x02 != 0)
	e.pins.Set(DB2, data&0x04 != 0)
	e.pins.Set(DB3, data&0x08 != 0)
	e.pins.Set(DB4, data&0x10 != 0)
	e.pins.Set(DB5, data&0x20 != 0)
	e.pins.Set(DB6, data&0x40 != 0)
	e.pins.Set(DB7, data&0x80 != 0)
	return e.pulseWrite()
}

func (e *Expander) write4Bits(data byte) error {
	e.pins.Set(DB4, data&0x01 != 0)
	e.pins.Set(DB5, data&0x02 != 0)
	e.pins.Set(DB6, data&0x04 != 0)
	e.pins.Set(DB7, data&0x08 != 0)
	return e.pulseWrite()
}

func (e *Expander) pulseWrite() error {
	err := e.w.WritePins(e.pins)
	if err != nil {
		return err
	}

	e.pins.Set(E, true)
	err = e.w.WritePins(e.pins)
	if err != nil {
		return err
	}

	e.pins.Set(E, false)
	return e.w.WritePins(e.pins)
}

func (e *Expander) writeByte(data byte) error {
	if e.eightBitMode {
		return e.write8Bits(data)
	}

	err := e.write4Bits(data >> 4)
	if err != nil {
		return err
	}

	return e.write4Bits(data)
}

func (e *Expander) readBits(pr ioexp.PinReader) (b byte, err error) {
	if e.eightBitMode {
		e.pins.Set(DB0, true)
		e.pins.Set(DB1, true)
		e.pins.Set(DB2, true)
		e.pins.Set(DB3, true)
	}
	e.pins.Set(DB4, true)
	e.pins.Set(DB5, true)
	e.pins.Set(DB6, true)
	e.pins.Set(DB7, true)
	err = e.w.WritePins(e.pins)
	if err != nil {
		return 0, err
	}

	e.pins.Set(E, true)
	err = e.w.WritePins(e.pins)
	if err != nil {
		return 0, err
	}

	pins, err := pr.ReadPins()
	if errors.Is(err, ioexp.ErrWriteOnly) {
		e.writeOnly = true
		// set E back to false
		e.pins.Set(E, false)
		if err := e.w.WritePins(e.pins); err != nil {
			return 0, err
		}
		return 0, err
	}
	if err != nil {
		return 0, err
	}

	e.pins.Set(E, false)
	err = e.w.WritePins(e.pins)

	if e.eightBitMode {
		return ioexp.PinByte(pins), err
	}
	return ioexp.PinByte(pins) & 0x0f, err
}

func (e *Expander) WriteByteIR(data byte) error {
	e.pins.Set(RS, false)
	e.pins.Set(RW, false)
	return e.writeByte(data)
}

func (e *Expander) WriteByte(data byte) error {
	e.pins.Set(RS, true)
	e.pins.Set(RW, false)
	return e.writeByte(data)
}

func (e *Expander) readByte() (b byte, err error) {
	pr, ok := e.w.(ioexp.PinReader)
	if !ok {
		e.writeOnly = true
		return 0, ioexp.ErrWriteOnly
	}
	if e.eightBitMode {
		return e.readBits(pr)
	}

	b1, err := e.readBits(pr)
	if err != nil {
		return 0, err
	}

	b, err = e.readBits(pr)
	if err != nil {
		return 0, err
	}

	return b1<<4 | b, nil
}

func (e *Expander) ReadByteIR() (byte, error) {
	if e.writeOnly {
		return 0, ioexp.ErrWriteOnly
	}

	e.pins.Set(RS, false)
	e.pins.Set(RW, true)

	return e.readByte()
}

func (e *Expander) ReadByte() (byte, error) {
	if e.writeOnly {
		return 0, ioexp.ErrWriteOnly
	}

	e.pins.Set(RS, true)
	e.pins.Set(RW, true)

	return e.readByte()
}
