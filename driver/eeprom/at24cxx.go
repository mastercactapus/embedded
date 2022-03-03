package eeprom

import (
	"fmt"
	"io"
	"time"
)

const timedWriteCycleDelay = 10 * time.Millisecond

type I2C interface {
	Tx(addr uint16, w, r []byte) error
}

type Device struct {
	i2c  I2C
	addr uint8

	maxLen  int
	pageLen int
	pos     int
	devPos  int

	wBuffer []byte
}

func NewAT24C01(i2c I2C, addr uint8) *Device {
	return &Device{
		i2c:     i2c,
		addr:    addr,
		pageLen: 8,
		maxLen:  128,
		devPos:  -1,
	}
}

func NewAT24C02(i2c I2C, addr uint8) *Device {
	return &Device{
		i2c:     i2c,
		addr:    addr,
		pageLen: 8,
		maxLen:  256,
		devPos:  -1,
	}
}

func NewAT24C04(i2c I2C, addr uint8) *Device {
	return &Device{
		i2c:     i2c,
		addr:    addr,
		pageLen: 16,
		maxLen:  512,
		devPos:  -1,
	}
}

func NewAT24C08(i2c I2C, addr uint8) *Device {
	return &Device{
		i2c:     i2c,
		addr:    addr,
		pageLen: 16,
		maxLen:  1024,
		devPos:  -1,
	}
}

func NewAT24C16(i2c I2C, addr uint8) *Device {
	return &Device{
		i2c:     i2c,
		addr:    addr,
		pageLen: 16,
		maxLen:  2048,
		devPos:  -1,
	}
}

func NewAT24C32(i2c I2C, addr uint8) *Device {
	return &Device{
		i2c:     i2c,
		addr:    addr,
		pageLen: 32,
		maxLen:  4096,
		devPos:  -1,
	}
}

func NewAT24C64(i2c I2C, addr uint8) *Device {
	return &Device{
		i2c:     i2c,
		addr:    addr,
		pageLen: 32,
		maxLen:  8192,
		devPos:  -1,
	}
}

var (
	_ = io.ReadWriteSeeker(&Device{})
	_ = io.ReaderAt(&Device{})
	_ = io.WriterAt(&Device{})
	_ = io.ByteWriter(&Device{})
	_ = io.ByteReader(&Device{})
)

func (d *Device) ReadAt(p []byte, off int64) (n int, err error) {
	_, err = d.Seek(off, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return d.Read(p)
}

func (d *Device) WriteAt(p []byte, off int64) (n int, err error) {
	_, err = d.Seek(off, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return d.Write(p)
}

func (d *Device) ReadByte() (c byte, err error) {
	buf := make([]byte, 1)
	_, err = d.Read(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (d *Device) Read(p []byte) (int, error) {
	if d.pos >= d.maxLen {
		return 0, io.EOF
	}
	if d.pos != d.devPos {
		// seek to correct location
		_, err := d._write(nil)
		if err != nil {
			return 0, err
		}
	}

	if len(p) > d.pageLen {
		p = p[:d.pageLen]
	}

	newPos := d.pos + len(p)
	if newPos > d.maxLen {
		p = p[:d.maxLen-d.pos]
	}

	err := d.i2c.Tx(uint16(d.addr), nil, p)
	if err != nil {
		d.devPos = -1
		return 0, err
	}

	d.pos = newPos
	d.devPos = d.pos
	return len(p), nil
}

func (d *Device) WriteByte(p byte) error {
	if d.pos >= d.maxLen {
		return io.ErrShortWrite
	}

	_, err := d._write([]byte{p})
	return err
}

func (d *Device) _write(p []byte) (n int, err error) {
	d.wBuffer = d.wBuffer[:0]

	addr := d.addr
	switch d.maxLen {
	case 128, 256:
		d.wBuffer = append(d.wBuffer, byte(d.pos))
	case 512:
		addr |= byte(d.pos>>8) & 0b1
		d.wBuffer = append(d.wBuffer, byte(d.pos))
	case 1024:
		addr |= byte(d.pos>>8) & 0b11
		d.wBuffer = append(d.wBuffer, byte(d.pos))
	case 2048:
		addr |= byte(d.pos>>8) & 0b111
		d.wBuffer = append(d.wBuffer, byte(d.pos))
	case 4096, 8192:
		d.wBuffer = append(d.wBuffer, byte(d.pos>>8), byte(d.pos))
	default:
		panic("unsupported device")
	}
	d.wBuffer = append(d.wBuffer, p...)
	err = d.i2c.Tx(uint16(addr), d.wBuffer, nil)
	if err != nil {
		d.devPos = -1
		return 0, err
	}

	d.pos += len(p)
	d.devPos = d.pos
	if len(p) > 0 {
		time.Sleep(timedWriteCycleDelay)
	}
	return len(p), nil
}

func (d *Device) Write(p []byte) (n int, err error) {
	// can't fit, return ErrShortWrite
	if d.pos+len(p) > d.maxLen {
		p = p[:d.maxLen-d.pos]
		if len(p) == 0 {
			return 0, io.ErrShortWrite
		}
		n, err = d.Write(p)
		if err != nil {
			return n, err
		}
		return n, io.ErrShortWrite
	}

	pageRem := d.pageLen - d.pos%d.pageLen
	if pageRem >= len(p) {
		// write current page if it fits
		return d._write(p)
	}

	// fill current page
	n, err = d._write(p[:pageRem])
	if err != nil {
		return n, err
	}
	p = p[pageRem:]

	// write full pages
	var _n int
	for len(p) >= d.pageLen {
		_n, err = d._write(p[:d.pageLen])
		n += _n
		if err != nil {
			return n, err
		}
		p = p[_n:]
	}

	if len(p) > 0 {
		_n, err = d._write(p)
		return n + _n, err
	}

	return n, nil
}

func (d *Device) Size() int {
	return d.maxLen
}

func (d *Device) Seek(offset int64, whence int) (int64, error) {
	if offset < 0 {
		return 0, fmt.Errorf("offset must be >= 0")
	}

	switch whence {
	case io.SeekStart:
		d.pos = int(offset)
	case io.SeekCurrent:
		d.pos += int(offset)
	case io.SeekEnd:
		d.pos = d.maxLen - int(offset)
	default:
		return 0, fmt.Errorf("invalid whence")
	}

	return int64(d.pos), nil
}
