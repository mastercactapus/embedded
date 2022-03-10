package mem

import (
	"fmt"
	"io"
	"time"

	"github.com/mastercactapus/embedded/serial"
)

// Pager is a memory device that can be used to read and write to pages of memory.
//
// Read and Write operations are performed by sending the memory start address.
//
// For devices that span multiple I2C addresses, use NewMultiPager to join multiple memory banks
// as descrete devices. See the AT24C04D, AT24C08D, and AT24C16D for examples.
type Pager struct {
	rw io.ReadWriter

	totalSize int
	pageSize  int
	addrSize  int
	devPos    int
	pos       int
	buf       []byte

	delay time.Duration
}

type PagerConfig struct {
	// PageSize is the size of a memory page in bytes.
	//
	// This controls the maximum chunk size and alignment when
	// performing writes.
	PageSize int

	// Capacity is the total size of the memory device in bytes.
	//
	// For example, the AT24C01 has a capacity of 128 bytes (1024 bits).
	Capacity int

	// AddressSize is the size of the address used to access the memory device.
	AddressSize int

	// WriteDelay is the time after writing to wait for the timed write cycle to complete.
	WriteDelay time.Duration
}

func NewPager(rw io.ReadWriter, cfg PagerConfig) *Pager {
	return &Pager{
		rw:        rw,
		buf:       make([]byte, cfg.PageSize+cfg.AddressSize),
		addrSize:  cfg.AddressSize,
		pageSize:  cfg.PageSize,
		totalSize: cfg.Capacity,
		devPos:    -1,
		delay:     cfg.WriteDelay,
	}
}

func (d *Pager) eof() bool {
	return d.pos >= d.totalSize
}

func (d *Pager) Read(p []byte) (n int, err error) {
	if d.eof() {
		return 0, io.EOF
	}

	if len(p) > d.remBytes() {
		p = p[:d.remBytes()]
	}

	if d.devPos != d.pos {
		switch d.addrSize {
		case 1:
			err = serial.Tx(d.rw, []byte{byte(d.pos)}, p)
		case 2:
			err = serial.Tx(d.rw, []byte{byte(d.pos >> 8), byte(d.pos & 0xff)}, p)
		default:
			panic("invalid address size")
		}
		if err != nil {
			return 0, err
		}
		d.incrPos(len(p))
		return len(p), err
	}

	n, err = d.rw.Read(p)
	d.incrPos(n)
	return n, err
}

func (d *Pager) remBytes() int {
	return d.totalSize - d.pos
}

func (d *Pager) remPageBytes() int {
	return d.pageSize - (d.pos % d.pageSize)
}

// pageBuf returns a buffer for the remaining page size + 1.
func (d *Pager) pageBuf() []byte {
	return d.buf[:d.remPageBytes()+1]
}

func (d *Pager) incrPos(n int) {
	d.pos += n
	if d.pos == d.totalSize {
		d.devPos = 0
	} else {
		d.devPos = d.pos
	}
}

func (d *Pager) ReadAt(p []byte, offset int64) (_ int, err error) {
	if offset != int64(d.pos) {
		_, err = d.Seek(offset, io.SeekStart)
		if err != nil {
			return 0, err
		}
	}

	return d.Read(p)
}

func (d *Pager) WriteAt(p []byte, offset int64) (_ int, err error) {
	if offset != int64(d.pos) {
		_, err = d.Seek(offset, io.SeekStart)
		if err != nil {
			return 0, err
		}
	}

	return d.Write(p)
}

func (d *Pager) Write(p []byte) (n int, err error) {
	if d.eof() {
		return 0, io.EOF
	}

	if len(p) > d.remBytes() {
		p = p[:d.remBytes()]
		n, err = d.Write(p)
		d.incrPos(n)
		if err != nil {
			return n, err
		}
		return n, io.ErrShortWrite
	}

	rem := p
	for len(rem) > 0 {
		pbuf := d.pageBuf()
		switch d.addrSize {
		case 1:
			pbuf[0] = byte(d.pos)
		case 2:
			pbuf[0] = byte(d.pos >> 8)
			pbuf[1] = byte(d.pos & 0xff)
		default:
			panic("invalid address size")
		}

		n, err = d.rw.Write(pbuf[:copy(pbuf[d.addrSize:], rem)+d.addrSize])
		d.incrPos(n - d.addrSize)
		if err != nil {
			// sent bytes -1 for the address byte
			return n + len(p) - len(rem) - d.addrSize, err
		}
		if d.delay > 0 {
			time.Sleep(d.delay)
		}

		rem = rem[n-1:]
	}

	return len(p), nil
}

func (d *Pager) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekEnd:
		return d.Seek(int64(d.totalSize)-offset, io.SeekStart)
	case io.SeekCurrent:
		return d.Seek(int64(d.pos)+offset, io.SeekStart)
	case io.SeekStart:
	default:
		return 0, fmt.Errorf("invalid whence")
	}
	if offset < 0 {
		return 0, fmt.Errorf("out of bounds")
	}

	d.pos = int(offset)

	return offset, nil
}
