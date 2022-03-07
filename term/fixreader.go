package term

import (
	"io"
	"time"
)

// fixReader fixes broken reader implementations.
type fixReader struct {
	io.Reader
}

func (f *fixReader) Read(p []byte) (n int, err error) {
	for n == 0 {
		n, err = f.Reader.Read(p)
		if err != nil {
			return n, err
		}
		if n == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}

	return n, err
}
