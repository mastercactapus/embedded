package ansi

import (
	"bytes"
)

func CutAfter(p []byte, sep byte) ([]byte, []byte) {
	idx := bytes.IndexByte(p, sep)
	if idx < 0 {
		return p, nil
	}

	return p[:idx+1], p[idx+1:]
}

func Cut(s string, sep byte) (string, string) {
	idx := bytes.IndexByte([]byte(s), sep)
	if idx < 0 {
		return s, ""
	}

	return s[:idx], s[idx+1:]
}
