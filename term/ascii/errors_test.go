package ascii

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorf(t *testing.T) {
	err := Errorf("foobar")
	assert.Equal(t, "foobar", err.Error())

	bazErr := errors.New("baz")
	err = Errorf("foobar: %w", bazErr)
	assert.Equal(t, "foobar: baz", err.Error())

	assert.True(t, errors.Is(err, bazErr))
}
