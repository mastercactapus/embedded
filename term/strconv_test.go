package term

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItoa(t *testing.T) {
	assert.Equal(t, "0", itoa(0))
	assert.Equal(t, "-1", itoa(-1))
	assert.Equal(t, "1", itoa(1))
	assert.Equal(t, "123", itoa(123))
	assert.Equal(t, "-123", itoa(-123))

	assert.Equal(t, "-123456", itoa(-123456))
	assert.Equal(t, "123456", itoa(123456))
}
