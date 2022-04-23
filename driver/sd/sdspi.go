package sd

import (
	"io"
	"time"

	"github.com/mastercactapus/embedded/driver"
	"github.com/mastercactapus/embedded/serial/spi"
	"github.com/mastercactapus/embedded/term/ascii"
)

type SPISD struct {
	spi spi.Controller
	cs  driver.Pin
	err error

	crc byte

	totalBlocks uint32
	block       uint32
	offset      uint32

	sdhc bool
}

func NewSPISD(spi spi.Controller, cs driver.Pin) *SPISD {
	cs.Output()
	cs.High()
	return &SPISD{spi: spi, cs: cs}
}

func (sd *SPISD) readErr() error {
	err := sd.err
	sd.err = nil
	return err
}

func (sd *SPISD) setCS(v bool) {
	if sd.err != nil {
		return
	}
	if v {
		sd.err = sd.cs.High()
	} else {
		sd.err = sd.cs.Low()
	}
}

func (sd *SPISD) w(v byte) {
	if sd.err != nil {
		return
	}

	sd.crc = CRC7(sd.crc, v)
	_, sd.err = sd.spi.ReadWriteByte(v)
}

func (sd *SPISD) r() (v byte) {
	if sd.err != nil {
		return 0
	}

	v, sd.err = sd.spi.ReadWriteByte(0xff)
	return v
}

// Init puts the card into SPI mode.
func (sd *SPISD) Init() error {
	sd.setCS(true)
	// need to pulse at least 74 cycles to initialize
	for i := 0; i < 10; i++ {
		sd.w(0xff)
	}
	sd.setCS(false)

	s := time.Now()
	for sd.r() != 0xff {
		if sd.err != nil {
			return sd.readErr()
		}
		if time.Since(s) > time.Second {
			return ascii.Errorf("sd: timeout waiting for init")
		}
	}

	for i := 25; i != 0; i-- {
		res, err := sd.writeCommand(0, 0)
		if err != nil {
			return err
		}
		if res == 0x01 {
			break
		}

		if i == 0 {
			return ascii.Errorf("sd: init failed CMD0: %#x", byte(res))
		}
	}

	res, err := sd.writeCommand(8, 0x1aa)
	if err != nil {
		return err
	}

	var arg uint32
	if res&0x04 == 0 {
		// high capacity
		sd.r()
		sd.r()
		sd.r()
		if sd.r() != 0xaa {
			return ascii.Errorf("sd: init failed CMD8: %#x", byte(res))
		}

		sd.sdhc = true
		arg = 0x40000000
	}

	for {
		res, err := sd.aCmd(41, arg)
		if err != nil {
			return err
		}
		if res == 0 {
			break
		}
	}

	_, err = sd.writeCommand(9, 0)
	if err != nil {
		return err
	}
	if err = sd.waitDataStart(); err != nil {
		return err
	}

	var csd CSD
	for i := range csd {
		csd[i] = sd.r()
	}
	sd.totalBlocks = csd.Blocks()

	// crc bytes
	sd.r()
	sd.r()

	return sd.readErr()
}

func (sd *SPISD) aCmd(cmd byte, arg uint32) (res R1, err error) {
	res, err = sd.writeCommand(55, 0)
	if err != nil {
		return
	}

	return sd.writeCommand(cmd, arg)
}

func (sd *SPISD) waitNotBusy() error {
	s := time.Now()
	for sd.r() != 0xff {
		if sd.err != nil {
			return sd.readErr()
		}
		if time.Since(s) > time.Second {
			return ascii.Errorf("sd: timeout waiting for not busy")
		}
	}
	return nil
}

func (sd *SPISD) waitDataStart() error {
	s := time.Now()
	for sd.r() != 0xfe {
		if sd.err != nil {
			return sd.readErr()
		}
		if time.Since(s) > time.Second {
			return ascii.Errorf("sd: timeout waiting for start")
		}
	}
	return nil
}

func (sd *SPISD) writeCommand(cmd byte, arg uint32) (res R1, err error) {
	sd.setCS(false)

	if err = sd.waitNotBusy(); err != nil {
		return
	}

	// clear bit 7
	cmd &= (1 << 7) ^ 0xFF
	// set bit 6
	cmd |= (1 << 6)

	sd.crc = 0
	sd.w(cmd)
	sd.w(byte(arg >> 24))
	sd.w(byte(arg >> 16))
	sd.w(byte(arg >> 8))
	sd.w(byte(arg))
	sd.w(sd.crc | 1)

	s := time.Now()
	for time.Since(s) < 5*time.Second {
		res = R1(sd.r())
		if sd.err != nil {
			return res, sd.readErr()
		}
		if res&0x80 != 0 {
			continue
		}

		return res, sd.readErr()
	}

	return res, ascii.Errorf("sd: timeout waiting for response")
}

func (sd *SPISD) Seek(offset int64, whence int) (pos int64, err error) {
	switch whence {
	case io.SeekStart:
		pos = offset
	case io.SeekCurrent:
		pos = int64(sd.block)*512 + int64(sd.offset) + offset
	case io.SeekEnd:
		pos = int64(sd.totalBlocks)*512 + offset
	}
	sd.block = uint32(pos / 512)
	sd.offset = uint32(pos % 512)

	return pos, nil
}

func (sd *SPISD) Read(p []byte) (n int, err error) {
	if !sd.sdhc {
		return 0, ascii.Errorf("sd: read not supported for non-SDHC cards")
	}

	if len(p) == 0 {
		return 0, nil
	}
	if sd.block >= sd.totalBlocks {
		return 0, io.EOF
	}

	var buf [512]byte
	err = sd.readBlock(sd.block, buf[:])
	if err != nil {
		return 0, err
	}

	n = copy(p, buf[sd.offset:])
	sd.advance(n)

	return n, sd.readErr()
}

func (sd *SPISD) advance(n int) (eof bool) {
	sd.offset += uint32(n)
	if sd.offset >= 512 {
		sd.block++
		sd.offset -= 512
	}
	return sd.block >= sd.totalBlocks
}

func (sd *SPISD) readBlock(block uint32, data []byte) error {
	if len(data) != 512 {
		return ascii.Errorf("sd: read block: invalid data length")
	}
	res, err := sd.writeCommand(17, block)
	defer sd.cs.Set(true)
	if err != nil {
		return err
	}
	if res != 0 {
		return ascii.Errorf("sd: read: %#x", byte(res))
	}

	if err = sd.waitDataStart(); err != nil {
		return err
	}

	_, err = spi.Read(sd.spi, 0xff, data)
	if err != nil {
		return err
	}

	// TODO: CRC
	sd.r()
	sd.r()

	sd.setCS(true)
	return sd.readErr()
}

func (sd *SPISD) writeBlock(block uint32, data []byte) error {
	if len(data) != 512 {
		return ascii.Errorf("sd: write block: invalid data length %d", len(data))
	}

	res, err := sd.writeCommand(24, block)
	defer sd.cs.Set(true)
	if err != nil {
		return ascii.Errorf("sd: write block: CMD24: %w", err)
	}
	if res != 0 {
		return ascii.Errorf("sd: write block: %#x", byte(res))
	}

	if err = sd.waitNotBusy(); err != nil {
		return err
	}

	// start token
	sd.w(0xfe)

	_, err = spi.Write(sd.spi, data)
	if err != nil {
		return err
	}

	sd.w(0xff)
	sd.w(0xff)
	if sd.err != nil {
		return sd.readErr()
	}

	sd.r() // TODO: check R1 response

	if err = sd.waitNotBusy(); err != nil {
		return err
	}

	sd.setCS(true)
	return sd.readErr()
}

func (sd *SPISD) Write(p []byte) (int, error) {
	if !sd.sdhc {
		return 0, ascii.Errorf("sd: read not supported for non-SDHC cards")
	}
	if len(p) == 0 {
		return 0, nil
	}
	if sd.block >= sd.totalBlocks {
		return 0, io.EOF
	}

	if sd.offset > 0 {
		var buf [512]byte
		err := sd.readBlock(sd.block, buf[:])
		if err != nil {
			return 0, ascii.Errorf("sd: read block: %w", err)
		}
		n := copy(buf[sd.offset:], p)
		err = sd.writeBlock(sd.block, buf[:])
		if err != nil {
			return 0, ascii.Errorf("sd: write block: %w", err)
		}
		sd.advance(n)

		if len(p) <= n {
			return len(p), nil
		}

		nx, err := sd.Write(p[n:])
		return n + nx, err
	}

	if len(p) < 512 {
		var buf [512]byte
		err := sd.readBlock(sd.block, buf[:])
		if err != nil {
			return 0, ascii.Errorf("sd: read block: %w", err)
		}
		n := copy(buf[:], p)
		err = sd.writeBlock(sd.block, buf[:])
		if err != nil {
			return 0, ascii.Errorf("sd: write block: %w", err)
		}
		sd.advance(n)
		return n, nil
	}

	if err := sd.writeBlock(sd.block, p[:512]); err != nil {
		return 0, ascii.Errorf("sd: write block: %w", err)
	}
	sd.advance(512)

	if len(p) <= 512 {
		return len(p), nil
	}

	n, err := sd.Write(p[512:])
	return 512 + n, err
}

var _ io.ReadWriteSeeker = (*SPISD)(nil)
