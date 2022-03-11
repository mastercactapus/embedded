package lcd

import (
	"errors"
	"time"

	"github.com/mastercactapus/embedded/driver/ioexp"
	"github.com/mastercactapus/embedded/serial/i2c"
)

type HD44780 struct {
	c         Controller
	writeOnly bool
	err       error
}

func NewHD44780I2C(bus i2c.Bus, addr uint16) *HD44780 {
	return NewHD44780(NewExpander(ioexp.NewPCF8574(bus, addr)))
}

func NewHD44780(c Controller) *HD44780 {
	return &HD44780{
		c: c,
	}
}

func (lcd *HD44780) Init() error {
	lcd.c.SetBacklight(true)
	time.Sleep(50 * time.Millisecond)
	if !lcd.c.IsEightBitMode() {
		// hack to get 4-bit mode from unknown state
		b := byte(0x03<<4) | 0x03
		lcd.writeIR(b, 5*time.Millisecond)
		b = 0x03<<4 | 0x02
		lcd.writeIR(b, 5*time.Millisecond)
	}
	lcd.SetFunction(lcd.c.IsEightBitMode(), true, false)
	lcd.SetDisplay(true, false, false)
	lcd.Clear()
	lcd.SetEntryMode(true, false)
	lcd.Home()

	return lcd.err
}

func (lcd *HD44780) SetBacklight(v bool) error {
	return lcd.c.SetBacklight(v)
}

func (lcd *HD44780) writeIR(b byte, dur time.Duration) error {
	if lcd.err != nil {
		return lcd.err
	}
	lcd.err = lcd.c.WriteByteIR(b)
	if lcd.err != nil {
		return lcd.err
	}
	return lcd.waitFor(dur)
}

func (lcd *HD44780) waitFor(dur time.Duration) error {
	if lcd.writeOnly {
		time.Sleep(dur)
		return nil
	}

	for {
		b, err := lcd.c.ReadByteIR()
		if errors.Is(err, ioexp.ErrWriteOnly) {
			lcd.writeOnly = true
			time.Sleep(dur)
			break
		}
		if err != nil {
			lcd.err = err
			return err
		}
		if b&(1<<7) == 0 {
			break
		}
	}

	return nil
}

// Clear clears the entire display and sets DDRAM address 0 in the address counter.
func (lcd *HD44780) Clear() error { return lcd.writeIR(0b_0000_0001, 2*time.Millisecond) }

// Home sets DDRAM address 0 in the address counter and returns the cursor to the home position.
func (lcd *HD44780) Home() error { return lcd.writeIR(0b_0000_0010, 1520*time.Microsecond) }

func (lcd *HD44780) SetEntryMode(increment, shiftDisplay bool) error {
	ins := byte(0b_0000_0100)
	if increment {
		ins |= 0b_0000_0010
	}
	if shiftDisplay {
		ins |= 0b_0000_0001
	}
	return lcd.writeIR(ins, 37*time.Microsecond)
}

func (lcd *HD44780) SetDisplay(displayOn, cursorOn, cursorBlink bool) error {
	ins := byte(0b_0000_1000)
	if displayOn {
		ins |= 0b_0000_0100
	}
	if cursorOn {
		ins |= 0b_0000_0010
	}
	if cursorBlink {
		ins |= 0b_0000_0001
	}
	return lcd.writeIR(ins, 37*time.Microsecond)
}

// SetFunction sets the number of display lines, character font, and data transmission length.
//
// If eightBitMode is false, 4-bit communication will be used.
// If twoLines is false, only one line will be displayed.
// If fiveByTen is false, the character font will be 5x8.
func (lcd *HD44780) SetFunction(eightBitMode, twoLines, fiveByTen bool) error {
	ins := byte(0b_0010_0000)
	if eightBitMode {
		ins |= 0b_0001_0000
	}
	if twoLines {
		ins |= 0b_0000_1000
	}
	if fiveByTen {
		ins |= 0b_0000_0100
	}
	return lcd.writeIR(ins, 37*time.Microsecond)
}

func (lcd *HD44780) SetCGRAMAddr(addr byte) error {
	if addr > 0x3F {
		return errors.New("addr must be 0 to 0x3F")
	}
	return lcd.writeIR(0b_0100_0000|addr, 37*time.Microsecond)
}

func (lcd *HD44780) SetDDRAMAddr(addr byte) error {
	if addr > 0x7F {
		return errors.New("invalid DDRAM address")
	}
	return lcd.writeIR(0b_1000_0000|addr, 37*time.Microsecond)
}

func (lcd *HD44780) SetCursor(col, line byte) error {
	if line > 1 {
		return errors.New("line must be 0 or 1")
	}
	if col > 15 {
		return errors.New("col must be 0 to 15")
	}
	return lcd.SetDDRAMAddr(col + (line * 0x40))
}

func (lcd *HD44780) WriteByte(b byte) error {
	defer lcd.waitFor(37 * time.Microsecond)
	return lcd.c.WriteByte(b)
}

func (lcd *HD44780) Write(p []byte) (int, error) {
	for i, b := range p {
		if err := lcd.WriteByte(b); err != nil {
			return i, err
		}
	}
	return len(p), nil
}
