package lcd

import (
	"errors"
	"time"

	"github.com/mastercactapus/embedded/driver"
	"github.com/mastercactapus/embedded/driver/ioexp"
	"github.com/mastercactapus/embedded/serial/i2c"
	"github.com/mastercactapus/embedded/term/ascii"
)

type HD44780 struct {
	c         Controller
	writeOnly bool
	err       error
	cfg       Config

	posX, posY int
}

type Config struct {
	Lines int
	Cols  int
	Font  Font
}

type Font int

const (
	Font5x8 Font = iota
	Font5x10
)

func NewHD44780I2C(bus i2c.Bus, addr uint16, cfg Config) (*HD44780, error) {
	pcf := ioexp.NewPCF8574(bus, addr)
	return NewHD44780(NewExpander(ExpanderConfig{
		RS:  pcf.BufferedPin(0),
		RW:  pcf.BufferedPin(1),
		E:   pcf.BufferedPin(2),
		BL:  pcf.BufferedPin(3),
		DB4: pcf.BufferedPin(4),
		DB5: pcf.BufferedPin(5),
		DB6: pcf.BufferedPin(6),
		DB7: pcf.BufferedPin(7),

		Flush:   pcf.Flush,
		Refresh: pcf.Refresh,
	}), cfg)
}

func NewHD44780(c Controller, cfg Config) (*HD44780, error) {
	lcd := &HD44780{
		c:   c,
		cfg: cfg,
	}

	if cfg.Lines > 4 {
		return nil, errors.New("lcd: only up t0 4 lines supported")
	}
	if cfg.Cols > 64 {
		return nil, errors.New("lcd: only up to 64 columns supported")
	}
	if cfg.Font != Font5x8 && cfg.Font != Font5x10 {
		return nil, errors.New("lcd: only 5x8 and 5x10 fonts supported")
	}

	if err := lcd.init(); err != nil {
		return nil, err
	}

	return lcd, nil
}

func (lcd *HD44780) init() error {
	lcd.err = nil
	lcd.c.SetBacklight(true)
	time.Sleep(50 * time.Millisecond)
	if !lcd.c.IsEightBitMode() {
		// hack to get 4-bit mode from unknown state
		b := byte(0x03<<4) | 0x03
		lcd.writeIR(b, 5*time.Millisecond)
		b = 0x03<<4 | 0x02
		lcd.writeIR(b, 5*time.Millisecond)
	}
	lcd.SetFunction(lcd.c.IsEightBitMode(), lcd.cfg.Lines > 1, lcd.cfg.Font == Font5x10)
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
		if errors.Is(err, driver.ErrNotSupported) {
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
func (lcd *HD44780) Home() error {
	lcd.posX, lcd.posY = 0, 0
	return lcd.writeIR(0b_0000_0010, 1520*time.Microsecond)
}

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

func (lcd *HD44780) _SetCGRAMAddr(addr byte) error {
	if addr > 0x3F {
		return errors.New("addr must be 0 to 0x3F")
	}
	return lcd.writeIR(0b_0100_0000|addr, 37*time.Microsecond)
}

func (lcd *HD44780) _SetDDRAMAddr(addr int) error {
	if addr >= 0x40+lcd.cfg.Cols*2 {
		return ascii.Errorf("invalid DDRAM address: %x", addr)
	}

	return lcd.writeIR(0b_1000_0000|byte(addr), 37*time.Microsecond)
}

// SetCursorXY sets the cursor position to the specified coordinates starting from top-left.
func (lcd *HD44780) SetCursorXY(x, y int) error {
	if y >= lcd.cfg.Lines {
		return errors.New("invalid line")
	}
	if x >= lcd.cfg.Cols {
		return errors.New("invalid column")
	}

	lcd.posX = x
	lcd.posY = y

	switch y {
	case 0:
		return lcd._SetDDRAMAddr(x)
	case 1:
		return lcd._SetDDRAMAddr(0x40 + x)
	case 2:
		return lcd._SetDDRAMAddr(x + lcd.cfg.Cols)
	case 3:
		return lcd._SetDDRAMAddr(0x40 + x + lcd.cfg.Cols)
	}
	panic("unreachable")
}

func (lcd *HD44780) incrX() error {
	lcd.posX++
	if lcd.posX < lcd.cfg.Cols {
		return nil
	}
	return lcd.incrY()
}

func (lcd *HD44780) incrY() error {
	lcd.posY++
	lcd.posX = 0
	if lcd.posY == lcd.cfg.Lines {
		lcd.posY = 0
		lcd.posX = 0
	}
	return lcd.SetCursorXY(lcd.posX, lcd.posY)
}

func (lcd *HD44780) WriteByte(b byte) error {
	defer lcd.waitFor(37 * time.Microsecond)
	if b == '\r' {
		lcd.posX = 0
		if err := lcd.SetCursorXY(0, lcd.posY); err != nil {
			return err
		}
		return nil
	}

	if b == '\n' {
		return lcd.incrY()
	}

	err := lcd.c.WriteByte(b)
	if err != nil {
		return err
	}

	return lcd.incrX()
}

func (lcd *HD44780) Write(p []byte) (int, error) {
	for i, b := range p {
		if err := lcd.WriteByte(b); err != nil {
			return i, err
		}
	}
	return len(p), nil
}
