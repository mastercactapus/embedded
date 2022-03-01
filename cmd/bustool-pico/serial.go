package main

import "io"

type fixReader struct {
	io.Reader
}

func (f *fixReader) Read(p []byte) (n int, err error) {
	for n == 0 {
		n, err = f.Reader.Read(p)
		if err != nil {
			return n, err
		}
	}

	return n, err
}
