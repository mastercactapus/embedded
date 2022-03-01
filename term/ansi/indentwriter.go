package ansi

import (
	"bytes"
	"io"
)

// IndentWriter is a Writer that indents all lines written to it.
type IndentWriter struct {
	w       io.Writer
	prefix  []byte
	partial bool
}

func NewIndentWriter(w io.Writer, indent string) *IndentWriter {
	return &IndentWriter{w: w, prefix: []byte(indent)}
}

func (w *IndentWriter) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	var buf bytes.Buffer

	if !w.partial {
		buf.Write(w.prefix)
	}

	var line, rem []byte
	rem = p
	for len(rem) > 0 {
		line, rem = CutAfter(rem, '\n')
		buf.Write(line)
		if len(rem) > 0 {
			buf.Write(w.prefix)
		}
	}
	w.partial = len(line) > 0 && line[len(line)-1] != '\n'

	n, err := w.w.Write(buf.Bytes())
	if err == nil {
		return len(p), nil
	}

	// TODO: calculate the number of bytes written
	return n, err
}
