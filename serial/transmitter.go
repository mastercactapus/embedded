package serial

import (
	"io"
)

type Transmitter interface {
	Tx(w, r []byte) error
}

// Tx can be used to transact a write and read together with a serial device.
func Tx(rw io.ReadWriter, w, r []byte) (err error) {
	if tr, ok := rw.(Transmitter); ok {
		return tr.Tx(w, r)
	}

	if len(w) > 0 {
		_, err = rw.Write(w)
		if err != nil {
			return err
		}
	}

	if len(r) > 0 {
		_, err = io.ReadFull(rw, r)
		if err != nil {
			return err
		}
	}

	return nil
}
