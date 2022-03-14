package ascii

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSprint(t *testing.T) {
	check := func(args ...interface{}) {
		assert.Equal(t, fmt.Sprint(args...), Sprint(args...))
	}

	check("foo")
	check("foo", 4)
	check("foo", []byte{1, 2, 3})
}

func TestSprintf(t *testing.T) {
	check := func(format string, args ...interface{}) {
		assert.Equal(t, fmt.Sprintf(format, args...), Sprintf(format, args...))
	}

	check("foo")
	check("foo bar%%")

	check("foo %s bar", "baz")
	check("foo %v bar", "baz")

	check("foo %c bar", 'h')
	check("foo %v bar", 'h')

	check("foo %x bar", 1)
	check("foo %x bar", "hi")
	check("foo %x bar", []byte("hi"))
	check("foo %#x bar", 1)
	check("foo %#02x bar", 1)
	check("foo %#04x bar", 1345)

	check("foo %d bar", "f")
	check("foo %s bar", 1)

	check("foo %d bar", 1)
	check("foo %d bar", 123)
	check("foo %-4d bar", 1)
	check("foo %-04d bar", 1)
	check("foo %4d bar", 1)
	check("foo %+4d bar", 1)
	check("foo %4d bar", -1)
	check("foo %+4d bar", -1)
	check("foo %04d bar", 1)
	check("foo %+04d bar", 1)

	check("foo %+04d bar", 1, 3)
}
