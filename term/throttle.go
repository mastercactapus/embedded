package term

import (
	"io"
	"time"
)

func NewThrottleWriter(w io.Writer, baud int) io.Writer {
	return &Throttle{Writer: w, delay: time.Second / time.Duration(baud)}
}

type Throttle struct {
	io.Writer
	delay time.Duration
}

func (t *Throttle) Write(p []byte) (int, error) {
	for i := range p {
		_, err := t.Writer.Write(p[i : i+1])
		if err != nil {
			return i, err
		}
		time.Sleep(t.delay)
	}

	return len(p), nil
}
