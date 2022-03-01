package term

import (
	"bufio"
)

type newliner struct {
	*bufio.Writer
	lastWByte byte
}

// Write writes to the underlying writer, but also ensures that newlines
// are always written with a preceding carriage return.
func (w *newliner) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	var err error
	for i, b := range p {
		if b == '\n' && w.lastWByte != '\r' {
			err = w.Writer.WriteByte('\r')
			if err != nil {
				return i, err
			}
		}
		w.lastWByte = b
		err = w.Writer.WriteByte(b)
		if err != nil {
			return i, err
		}
	}

	return len(p), err
}
