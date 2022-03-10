package onewire

import "encoding/binary"

type Address uint64

func (a Address) CRC() byte { return byte(a) }

// TODO: decode and provid family info
func (a Address) Family() byte { return byte(a & 0xff) }

func (a Address) Serial() uint64 {
	return uint64(a>>8) & 0xffffffffffff
}

func (a Address) Valid() bool {
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], uint64(a))

	return a.CRC() == CRC8(data[:7])
}
