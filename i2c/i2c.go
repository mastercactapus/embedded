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

func (i2c *I2C) WriteTo(p []byte, addr byte) (int, error) {
	i2c.Start()
	defer i2c.Stop()

	if err := i2c.WriteByte(addr << 1); err != nil {
		return 0, err
	}

	return i2c.Write(p)
}

func (i2c *I2C) Write(p []byte) (int, error) {
	for i, b := range p {
		if err := i2c.WriteByte(b); err != nil {
			return i, err
		}
	}
	return len(p), nil
}

func (i2c *I2C) ReadFrom(p []byte, addr byte) (int, error) {
	i2c.Start()
	defer i2c.Stop()

	if err := i2c.WriteByte((addr << 1) | 1); err != nil {
		return 0, err
	}

	return i2c.Read(p)
}

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

var ErrNack = errors.New("i2c: NACK")

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

func (i2c *I2C) Stop() {
	i2c.sda.OutputLow()
	wait()
	i2c.clockUp()
	wait()
	i2c.sda.PullupHigh()
	wait()
}
