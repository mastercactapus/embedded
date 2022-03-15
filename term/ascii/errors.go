package ascii

import (
	"bytes"
	"errors"
)

func Errorf(format string, a ...interface{}) error {
	var buf bytes.Buffer
	buf.Grow(len(format))
	f := &fPrinter{w: &buf, format: format, args: a, isErrorf: true}
	f.printf(format, a)
	if f.err != nil {
		// should not be possible
		panic(f.err)
	}
	if f.wrappedErr != nil {
		return &wrapError{buf.String(), f.wrappedErr}
	}

	return errors.New(buf.String())
}

type wrapError struct {
	msg string
	err error
}

func (e *wrapError) Error() string {
	return e.msg
}

func (e *wrapError) Unwrap() error {
	return e.err
}
