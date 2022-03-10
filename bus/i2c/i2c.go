package i2c

import (
	"errors"
)

type Controller interface {
	Start()
	Stop()
	WriteBit(bool)
	ReadBit() bool
}

type BaudRateController interface {
	SetBaudRate(baudrate uint32) error
}

type Pin interface {
	PullupHigh()
	OutputLow()

	Get() bool
}

type I2C struct {
	Controller
}

func New(c Controller) *I2C {
	return &I2C{Controller: c}
}

var (
	ErrBadAddr     = errors.New("i2c: bad address")
	ErrNack        = errors.New("i2c: NACK")
	ErrUnsupported = errors.New("i2c: unsupported")
)

func (i2c *I2C) SetBaudrate(baudrate uint32) error {
	if bc, ok := i2c.Controller.(BaudRateController); ok {
		return bc.SetBaudRate(baudrate)
	}

	return ErrUnsupported
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
func (i2c *I2C) Read(p []byte) (n int, err error) {
	for i := range p {
		if p[i], err = i2c._ReadByte(i == len(p)-1); err != nil {
			return i, err
		}
	}
	return len(p), nil
}

// ReadByte reads a single byte directly from the bus, without any address or start/stop condition.
func (i2c *I2C) ReadByte() (byte, error) {
	return i2c._ReadByte(true)
}

func (i2c *I2C) _ReadByte(nak bool) (byte, error) {
	var b byte
	for i := 0; i < 8; i++ {
		if i2c.ReadBit() {
			b |= 1 << (7 - i)
		}
	}

	i2c.WriteBit(nak)

	return b, nil
}

// WriteByte writes a single byte to the bus, without any address or start/stop condition.
func (i2c *I2C) WriteByte(b byte) error {
	for i := 0; i < 8; i++ {
		i2c.WriteBit(((b >> (7 - i)) & 1) == 1)
	}

	if i2c.ReadBit() {
		return ErrNack
	}

	return nil
}
