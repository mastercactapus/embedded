package mem_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/mastercactapus/embedded/driver/mem"
	"github.com/stretchr/testify/require"
)

type rws struct {
	io.ReadSeeker
}

func (rws) Write([]byte) (int, error) { return 0, nil }

func TestJoin(t *testing.T) {
	a := rws{bytes.NewReader([]byte{1, 2, 3})}
	b := rws{bytes.NewReader([]byte{4, 5, 6, 7})}
	c := rws{bytes.NewReader([]byte{8, 9, 10})}

	mp := mem.Join(a, b, c)

	data, err := io.ReadAll(mp)
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, data)

	_, err = mp.Seek(0, 0)
	require.NoError(t, err)

	data, err = io.ReadAll(mp)
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, data)

	mp.Seek(3, 0)
	data, err = io.ReadAll(mp)
	require.NoError(t, err)
	require.Equal(t, []byte{4, 5, 6, 7, 8, 9, 10}, data)
}
