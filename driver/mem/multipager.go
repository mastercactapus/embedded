package mem

import (
	"fmt"
	"io"
)

// MultiPager is a memory device that can be composed of multiple memory devices.
type MultiPager struct {
	rws []io.ReadWriteSeeker

	devMax     []int
	devLastPos []int

	devPos int
	devIdx int

	pos    int
	maxLen int
}

// Join will join the given memory devices into a single device.
func Join(rws ...io.ReadWriteSeeker) io.ReadWriteSeeker {
	if len(rws) == 1 {
		return rws[0]
	}
	m := &MultiPager{rws: rws}
	for _, rw := range rws {
		l, err := rw.Seek(0, io.SeekEnd)
		if err != nil {
			panic(err)
		}
		m.maxLen += int(l)
		m.devMax = append(m.devMax, int(l))
		_, err = rw.Seek(0, 0)
		if err != nil {
			panic(err)
		}
	}

	return m
}

func (m *MultiPager) eof() bool {
	return m.devIdx == len(m.rws)
}

func (m *MultiPager) needSeek() bool {
	return m.devPos != m.devLastPos[m.devIdx]
}

func (m *MultiPager) remBytes() int {
	return m.maxLen - m.pos
}

func (m *MultiPager) remDevBytes() int {
	return m.devMax[m.devIdx] - m.devLastPos[m.devIdx]
}
func (m *MultiPager) dev() io.ReadWriteSeeker { return m.rws[m.devIdx] }

func (m *MultiPager) incrPos(n int) {
	m.pos += n
	m.devPos += n
	m.devLastPos[m.devIdx] = m.devPos
	if n < m.remDevBytes() {
		return
	}

	if n != m.remDevBytes() {
		panic("unaligned position")
	}

	m.devIdx++
	m.devPos = 0
}

func (m *MultiPager) Read(p []byte) (n int, err error) {
	if m.eof() {
		return 0, io.EOF
	}

	if len(p) > m.remDevBytes() {
		return m.Read(p[:m.remDevBytes()])
	}

	if m.needSeek() {
		if ra, ok := m.dev().(io.ReaderAt); ok {
			n, err = ra.ReadAt(p, int64(m.devPos))
			m.incrPos(n)
			return n, err
		}

		np, err := m.dev().Seek(int64(m.devPos), 0)
		m.devLastPos[m.devIdx] = int(np)
		if err != nil {
			return 0, err
		}
	}

	n, err = m.dev().Read(p)
	m.incrPos(n)
	return n, err
}

func (m *MultiPager) ReadAt(p []byte, offset int64) (_ int, err error) {
	if offset != int64(m.pos) {
		_, err = m.Seek(offset, io.SeekStart)
		if err != nil {
			return 0, err
		}
	}

	return m.Read(p)
}

func (m *MultiPager) WriteAt(p []byte, offset int64) (_ int, err error) {
	if offset != int64(m.pos) {
		_, err = m.Seek(offset, io.SeekStart)
		if err != nil {
			return 0, err
		}
	}

	return m.Write(p)
}

func (m *MultiPager) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekEnd:
		return m.Seek(int64(m.maxLen)-offset, io.SeekStart)
	case io.SeekCurrent:
		return m.Seek(int64(m.pos)+offset, io.SeekStart)
	case io.SeekStart:
	default:
		return 0, fmt.Errorf("invalid whence")
	}
	if offset < 0 {
		return 0, fmt.Errorf("out of bounds")
	}

	if offset >= int64(m.maxLen) {
		m.devIdx = len(m.rws)
		m.pos = int(offset)
		m.devPos = 0
		return offset, nil
	}

	m.pos = int(offset)

	m.devIdx = 0
	for offset >= int64(m.devMax[m.devIdx]) {
		offset -= int64(m.devMax[m.devIdx])
		m.devIdx++
	}

	m.devPos = int(offset)

	return int64(m.pos), nil
}

func (m *MultiPager) Write(p []byte) (n int, err error) {
	if m.eof() {
		return 0, io.EOF
	}

	if len(p) > m.remBytes() {
		p = p[:m.remBytes()]
		n, err = m.Write(p)
		if err != nil {
			return n, err
		}
		return n, io.ErrShortWrite
	}

	buf := p
	for len(buf) > m.remDevBytes() {
		n, err = m.Write(buf[:m.remDevBytes()])
		if err != nil {
			return n + len(p) - len(buf), err
		}
		buf = buf[m.remDevBytes():]
	}

	if m.needSeek() {
		if wa, ok := m.dev().(io.WriterAt); ok {
			n, err = wa.WriteAt(buf, int64(m.devPos))
			m.incrPos(n)
			return n, err
		}

		_, err = m.dev().Seek(int64(m.devPos), 0)
		if err != nil {
			return 0, err
		}
	}

	n, err = m.dev().Write(buf)
	m.incrPos(n)
	return n + len(p) - len(buf), err
}
