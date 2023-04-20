package at_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/mastercactapus/embedded/at"
	"github.com/stretchr/testify/assert"
)

func TestServ(t *testing.T) {
	var buf bytes.Buffer
	s := at.NewServer(strings.NewReader("AT+TEST\r\n"), &buf)

	// Register a handler for the "ECHO" command
	s.HandleFunc("AT+TEST", func(c at.Cmd) at.Response {
		var resp at.Response
		resp.OK = true
		resp.SetValue("", "BAR")
		return resp
	})

	// Run the server
	err := s.Serve()
	assert.ErrorIs(t, err, io.ErrUnexpectedEOF)

	assert.Equal(t, "+TEST: =BAR\r\nOK\r\n", buf.String())
}
