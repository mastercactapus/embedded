package at_test

import (
	"testing"

	"github.com/mastercactapus/embedded/at"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	var resp at.Response

	resp.SetValue("foo", "bar")
	assert.Equal(t, "bar", resp.Value("foo"))

	resp.SetValue("", "1")
	assert.Equal(t, "1", resp.Value(""))
}
