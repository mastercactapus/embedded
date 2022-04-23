package sd_test

import (
	"testing"

	"github.com/mastercactapus/embedded/driver/sd"
	"github.com/stretchr/testify/assert"
)

func TestCRC7(t *testing.T) {
	// test with reset command
	var crc byte
	crc = sd.CRC7(crc, 0x40)
	crc = sd.CRC7(crc, 0x00)
	crc = sd.CRC7(crc, 0x00)
	crc = sd.CRC7(crc, 0x00)
	crc = sd.CRC7(crc, 0x00)
	crc |= 1 // set stop bit

	assert.Equal(t, byte(0x95), crc)
}
