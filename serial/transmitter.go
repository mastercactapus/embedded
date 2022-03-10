package serial

import "io"

type Transmitter interface {
	Tx(w, r []byte) error
}

func Tx(rw io.ReadWriter, w, r []byte) (err error) {
	if tr, ok := rw.(Transmitter); ok {
		return tr.Tx(w, r)
	}

	_, err = rw.Write(w)
	if err != nil {
		return err
	}

	_, err = rw.Read(r)
	return err
}
