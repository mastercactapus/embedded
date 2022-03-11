package mem_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/mastercactapus/embedded/driver/mem"
	"github.com/stretchr/testify/require"
)

func TestJoin(t *testing.T) {
	var a, b, c dev

	a.Write([]byte{1, 2, 3})
	b.Write([]byte{4, 5, 6})
	c.Write([]byte{7, 8, 9})

	mp := mem.Join(&a, &b, &c)

	data, err := io.ReadAll(mp)
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, data)

	_, err = mp.Seek(0, 0)
	require.NoError(t, err)

	data, err = io.ReadAll(mp)
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, data)

	mp.Seek(3, 0)
	data, err = io.ReadAll(mp)
	require.NoError(t, err)
	require.Equal(t, []byte{4, 5, 6, 7, 8, 9}, data)

	pos, err := mp.Seek(3, 0)
	require.NoError(t, err)
	require.Equal(t, int64(3), pos)
	n, err := mp.Write([]byte{10, 11, 12, 13})
	require.NoError(t, err)
	require.Equal(t, 4, n)

	pos, err = mp.Seek(2, 0)
	require.NoError(t, err)
	require.Equal(t, int64(2), pos)
	n, err = mp.Write([]byte{14, 15, 16})
	require.NoError(t, err)
	require.Equal(t, 3, n)

	_, err = mp.Seek(0, 0)
	require.NoError(t, err)

	data, err = io.ReadAll(mp)
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 14, 15, 16, 12, 13, 8, 9}, data)
}

type dev struct {
	buf [3]byte
	pos int
}

func (d *dev) Write(p []byte) (int, error) {
	n := copy(d.buf[d.pos:], p)
	d.pos += n
	if n < len(p) {
		return n, fmt.Errorf("tried to write %d bytes, but only wrote %d", len(p), n)
	}

	return n, nil
}

func (d *dev) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		d.pos = int(offset)
	case 1:
		d.pos += int(offset)
	case 2:
		d.pos = len(d.buf) - int(offset)
	}

	return int64(d.pos), nil
}

func (d *dev) Read(p []byte) (int, error) {
	if d.pos >= len(d.buf) {
		return 0, io.EOF
	}

	n := copy(p, d.buf[d.pos:])
	d.pos += n
	if n == 0 {
		return 0, io.EOF
	}

	return n, nil
}
