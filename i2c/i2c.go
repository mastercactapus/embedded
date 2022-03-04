package i2c

import (
	"errors"
)

type Pin interface {
	PullupHigh()
	OutputLow()

	Get() bool
}

type I2C struct {
	sda, scl Pin

	waitQtrN int
}

func New() *I2C {
	return &I2C{}
}

type Config struct {
	SDA, SCL Pin
	Baudrate uint32
}

var (
	ErrBadAddr = errors.New("i2c: bad address")
	ErrNack    = errors.New("i2c: NACK")
)

func (i2c *I2C) Configure(config Config) error {
	if config.Baudrate == 0 {
		config.Baudrate = 100e3
	}

	i2c.sda = config.SDA
	i2c.scl = config.SCL
	if i2c.sda == nil || i2c.scl == nil {
		return errors.New("i2c: pins not configured")
	}

	i2c.sda.PullupHigh()
	i2c.scl.PullupHigh()

	i2c.SetBaudrate(config.Baudrate)
	return nil
}

func (i2c *I2C) SetBaudrate(baudrate uint32) {
	// TODO
}

// TODO: set timeout based on baud
func (i2c *I2C) clockUp() {
	i2c.scl.PullupHigh()
	for !i2c.scl.Get() {
		// clock stretching
	}
}

const (
	Min7BitAddr = 0x08
	Max7BitAddr = 0x77

	Min10BitAddr = 0x7800
	Max10BitAddr = 0x7bff

	modeRead  = 1
	modeWrite = 0
)

const _10BitAddr = 0b11110_0000000000

func (i2c *I2C) writeLongAddress(addr uint16, mode byte) error {
	// first byte as write
	err := i2c.WriteByte(byte(addr>>7) & ^byte(1))
	if err != nil {
		return err
	}
	err = i2c.WriteByte(byte(addr & 0xff))
	if err != nil {
		return err
	}
	if mode == modeWrite {
		return nil
	}

	// read mode, retransmit first byte after start
	// with read flag set
	i2c.Start()
	err = i2c.WriteByte(byte(addr>>7) | 1)
	if err != nil {
		return err
	}

	return nil
}

func (i2c *I2C) writeAddress(addr uint16, mode byte) error {
	if addr < Min7BitAddr {
		// 7 not allowed
		// special cases, not suitible for normal addressing
		return ErrBadAddr
	}

	if addr > Max7BitAddr && addr < Min10BitAddr {
		// outside of 10-bit space
		return ErrBadAddr
	}

	if addr > Max10BitAddr {
		return ErrBadAddr
	}

	if addr > Max7BitAddr {
		return i2c.writeLongAddress(addr, mode)
	}

	return i2c.WriteByte(byte(addr<<1) | mode)
}

func (i2c *I2C) WriteByteTo(p byte, addr uint16) error {
	i2c.Start()
	defer i2c.Stop()

	if err := i2c.writeAddress(addr, modeWrite); err != nil {
		return err
	}

	return i2c.WriteByte(p)
}

func (i2c *I2C) WriteTo(p []byte, addr uint16) (int, error) {
	i2c.Start()
	defer i2c.Stop()

	if err := i2c.writeAddress(addr, modeWrite); err != nil {
		return 0, err
	}

	return i2c.Write(p)
}

func (i2c *I2C) ReadFrom(p []byte, addr uint16) (int, error) {
	i2c.Start()
	defer i2c.Stop()

	if err := i2c.writeAddress(addr, modeRead); err != nil {
		return 0, err
	}

	return i2c.Read(p)
}

func (i2c *I2C) ReadByteFrom(addr uint16) (byte, error) {
	i2c.Start()
	defer i2c.Stop()

	if err := i2c.writeAddress(addr, modeRead); err != nil {
		return 0, err
	}

	return i2c.ReadByte()
}

// Write will write directly to the bus, without any address or start/stop condition.
func (i2c *I2C) Write(p []byte) (int, error) {
	for i, b := range p {
		if err := i2c.WriteByte(b); err != nil {
			return i, err
		}
	}
	return len(p), nil
}

// Read will read directly from the bus, without any address or start/stop condition.
func (i2c *I2C) Read(p []byte) (int, error) {
	buf := p
	for len(buf) > 0 {
		for i := 0; i < 8; i++ {
			if i2c.readBit() {
				buf[0] = buf[0] | (1 << (7 - i))
			}
		}

		i2c.writeBit(len(buf) == 1)
		buf = buf[1:]
	}

	return len(p), nil
}

// ReadByte reads a single byte directly from the bus, without any address or start/stop condition.
func (i2c *I2C) ReadByte() (byte, error) {
	var b byte
	for i := 0; i < 8; i++ {
		if i2c.readBit() {
			b |= 1 << (7 - i)
		}
	}

	i2c.writeBit(true)

	wait()

	return b, nil
}

// WriteByte writes a single byte to the bus, without any address or start/stop condition.
func (i2c *I2C) WriteByte(b byte) error {
	for i := 0; i < 8; i++ {
		i2c.writeBit(((b >> (7 - i)) & 1) == 1)
	}

	if i2c.readBit() {
		return ErrNack
	}

	wait()

	return nil
}

// Start will send a start condition on the bus.
func (i2c *I2C) Start() {
	i2c.clockUp()
	wait()
	i2c.sda.OutputLow()
	wait()
	i2c.scl.OutputLow()
	wait()
}

func (i2c *I2C) writeBit(v bool) {
	if v {
		i2c.sda.PullupHigh()
	} else {
		i2c.sda.OutputLow()
	}
	wait()
	i2c.clockUp()
	wait()
	i2c.scl.OutputLow()
	wait()
}

func (i2c *I2C) readBit() (value bool) {
	i2c.sda.PullupHigh()
	wait()
	i2c.clockUp()
	wait()
	value = i2c.sda.Get()
	wait()
	i2c.scl.OutputLow()
	wait()
	return value
}

// Stop will send a stop condition on the bus.
func (i2c *I2C) Stop() {
	i2c.sda.OutputLow()
	wait()
	i2c.clockUp()
	wait()
	i2c.sda.PullupHigh()
	wait()
}
