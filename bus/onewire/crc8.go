package onewire

func CRC8(data []byte) byte {
	crc := byte(0)
	for _, b := range data {
		crc ^= b
		for i := 0; i < 8; i++ {
			if crc&0x01 != 0 {
				crc = (crc >> 1) ^ 0x8c
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}
