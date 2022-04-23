package sd

func CRC7(in, v byte) byte {
	in ^= v
	for i := 0; i < 8; i++ {
		if in&0x80 != 0 {
			in ^= 0x89
		}
		in <<= 1
	}
	return in
}
