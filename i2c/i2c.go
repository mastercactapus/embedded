package i2c

import (
	"errors"
	"time"
)

type Pin interface {
	PullupHigh()
	OutputLow()

	Get() bool
}

type I2C struct {
	sda, scl Pin
	waitDur  time.Duration
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

	// i2c.sda.PullupHigh()
	// i2c.scl.PullupHigh()

	i2c.SetBaudrate(config.Baudrate)
	return nil
}

func (i2c *I2C) SetBaudrate(baudrate uint32) {
	i2c.waitDur = time.Second / time.Duration(baudrate) / 4
}

func (i2c *I2C) waitHalf() {
	time.Sleep(i2c.waitDur)
	time.Sleep(i2c.waitDur)
}

func (i2c *I2C) waitQtr() {
	time.Sleep(i2c.waitDur)
}

func (i2c *I2C) Start() {
	// doubles as a STOP
	// if SDA was low
	i2c.clockUp()
	i2c.waitHalf()
	i2c.sda.PullupHigh()
	i2c.waitHalf()
	i2c.waitHalf()
	i2c.sda.OutputLow()
	i2c.waitHalf()
}

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

	return b, nil
}

func (i2c *I2C) readBit() bool {
	i2c.scl.OutputLow()
	i2c.waitHalf()
	i2c.sda.PullupHigh()
	i2c.clockUp()
	i2c.waitQtr()
	defer i2c.waitQtr()
	return i2c.sda.Get()
}

func (i2c *I2C) writeBit(v bool) {
	i2c.scl.OutputLow()
	if v {
		i2c.sda.PullupHigh()
	} else {
		i2c.sda.OutputLow()
	}
	i2c.waitHalf()
	i2c.clockUp()
	i2c.waitHalf()
}

func (i2c *I2C) WriteByte(b byte) error {
	for i := 0; i < 8; i++ {
		i2c.writeBit(((b >> (7 - i)) & 1) == 1)
	}

	if i2c.readBit() {
		return ErrNack
	}

	return nil
}

var ErrNack = errors.New("i2c: NACK")

func (i2c *I2C) Stop() {
	i2c.sda.OutputLow()
	i2c.scl.OutputLow()
	i2c.waitHalf()
	i2c.clockUp()
	i2c.waitHalf()
	i2c.sda.PullupHigh()
	i2c.waitHalf()
}
