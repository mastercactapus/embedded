package at

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeString(t *testing.T) {
	check := func(in string) {
		t.Helper()
		assert.Equal(t, in, UnescapeString(EscapeString(in)))
	}

	check("")
	check("foo")
	check("foo\nbar")
	check("foo\rbar")
	check("foo\\bar")
	check("foo\\bar\r\nbaz")
	check("foo\\ebar\r\nbaz\\qux")
}
