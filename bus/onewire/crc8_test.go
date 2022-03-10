package onewire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRC8(t *testing.T) {
	sum := CRC8([]byte{0x28, 0xa1, 0x72, 0x02, 0x00, 0x00, 0x00})
	assert.Equal(t, byte(0x9c), sum)
}
