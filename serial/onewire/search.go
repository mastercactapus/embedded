package onewire

import (
	"encoding/binary"

	"github.com/mastercactapus/embedded/term/ascii"
)

type searchState struct {
	alarm bool
	found []Address
	err   error
}

func (ow *OneWire) searchStart(s *searchState, rom uint64, startN int) {
	if s.err != nil {
		return
	}

	if !ow.Reset() {
		s.err = ErrNoDevice
		return
	}

	if s.alarm {
		s.err = ow.WriteByte(0xEC)
	} else {
		s.err = ow.WriteByte(0xF0)
	}

	ow.search(s, rom, startN, 0)
}

func (ow *OneWire) search(s *searchState, rom uint64, startN, n int) {
	if s.err != nil {
		return
	}
	if n == 64 {
		var tmp [8]byte
		binary.LittleEndian.PutUint64(tmp[:], rom)
		addr := Address(binary.BigEndian.Uint64(tmp[:]))
		if !addr.Valid() {
			s.err = ascii.Errorf("search [n=%d]: %x: %w", n, rom, ErrBadChecksum)
			return
		}

		s.found = append(s.found, addr)
		return
	}

	id, idC := ow.ReadBit(), ow.ReadBit()
	switch {
	case id && idC:
		s.err = ascii.Errorf("search [n=%d]: %w", n, ErrNoDevice)
		return
	case !id && !idC:
		if n < startN {
			ow.WriteBit(rom&(1<<uint(n)) != 0)
			ow.search(s, rom, startN, n+1)
			return
		}

		// search zero path first
		ow.WriteBit(false)
		ow.search(s, rom, startN, n+1)

		// search one path
		ow.searchStart(s, rom|(1<<uint(n)), n+1)
	case id:
		ow.WriteBit(true)
		ow.search(s, rom|(1<<uint(n)), startN, n+1)
	default:
		ow.WriteBit(false)
		ow.search(s, rom, startN, n+1)
	}
}
