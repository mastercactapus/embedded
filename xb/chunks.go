package xb

import (
	"bufio"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

type Header struct {
	Magic     byte
	Type      byte
	Length    uint32
	HeaderCRC uint32
	DataCRC   uint32
}

func (h *Header) CRC() uint32 {
	var crcPart [6]byte
	crcPart[0] = h.Magic
	crcPart[1] = h.Type
	binary.LittleEndian.PutUint32(crcPart[2:], h.Length)
	return crc32.ChecksumIEEE(crcPart[:])
}

func WriteChunk(w io.Writer, typeCode byte, data []byte) error {
	h := Header{
		Magic:  0xfe,
		Type:   typeCode,
		Length: uint32(len(data)),
	}
	h.HeaderCRC = h.CRC()
	h.DataCRC = crc32.ChecksumIEEE(data)

	err := binary.Write(w, binary.LittleEndian, h)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	if err != nil {
		return err
	}

	return err
}

func ReadChunk(r *bufio.Reader) (byte, []byte, error) {
	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, nil, err
		}
		if b == 0xfe {
			r.UnreadByte()
			break
		}
	}

	var h Header
	err := binary.Read(r, binary.LittleEndian, &h)
	if err != nil {
		return 0, nil, err
	}

	if h.HeaderCRC != h.CRC() {
		return 0, nil, errors.New("xb: chunk header CRC mismatch")
	}

	data := make([]byte, h.Length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return 0, nil, err
	}

	if h.DataCRC != crc32.ChecksumIEEE(data) {
		return 0, nil, errors.New("xb: chunk data CRC mismatch")
	}

	return h.Type, data, nil
}
