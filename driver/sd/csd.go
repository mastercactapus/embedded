package sd

type CSD [16]byte

// Version will return 0 for version 1.00 and
// 1 for version 2.00.
func (c CSD) Version() uint8 {
	return c[0] >> 6
}

func (c CSD) Blocks() uint32 {
	switch c.Version() {
	// case 0:
	case 1:
		return (uint32(c[7]>>2)<<16 | uint32(c[8])<<8 | uint32(c[9]) + 1) << 10
	}

	return 0
}
