package at

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"github.com/mastercactapus/embedded/term/ascii"
)

// Client is an AT command client.
type Client struct {
	rw io.ReadWriter
	s  *bufio.Scanner
	mx sync.Mutex
}

// NewClient creates a new AT command client.
func NewClient(rw io.ReadWriter) (*Client, error) {
	s := bufio.NewScanner(rw)
	s.Split(bufio.ScanLines)
	c := &Client{rw: rw, s: s}

	// Wait for the modem to send "HELLO".
	_, err := io.WriteString(c.rw, "\r\nSTART\r\n")
	if err != nil {
		return nil, err
	}
	for s.Scan() {
		if s.Text() == "HELLO" {
			return c, nil
		}
	}

	return nil, s.Err()
}

// Set sets the value of a parameter.
//
// The following example will result in sending "AT+FOO=bar\r\n"
// to the modem:
//
//	data, err := c.Set("foo", "bar")
func (c *Client) Set(name string, params ...string) (*Response, error) {
	name = strings.ToUpper(name)
	name = EscapeString(name, '=', '?')

	escaped := make([]string, 0, len(params))
	for _, param := range params {
		escaped = append(escaped, EscapeString(param, ','))
	}

	return c._Command(name, ascii.Sprintf("AT+%s=%s", name, strings.Join(escaped, ",")))
}

// Query queries the value of a parameter.
//
// The following example will result in sending "AT+FOO?\r\n"
// to the modem:
//
//	data, err := c.Query("foo")
func (c *Client) Query(name string) (*Response, error) {
	name = strings.ToUpper(name)
	name = EscapeString(name, '=', '?')

	return c._Command(name, ascii.Sprintf("AT+%s?", name))
}

// Execute executes a command.
//
// The following example will result in sending "AT+FOO\r\n"
// to the modem:
//
//	data, err := c.Execute("foo")
func (c *Client) Execute(name string) (*Response, error) {
	name = strings.ToUpper(name)
	name = EscapeString(name, '=', '?')

	return c._Command(name, ascii.Sprintf("AT+%s", name))
}

var ErrHello = ascii.Errorf("at: HELLO received")

func (c *Client) _Command(name, line string) (*Response, error) {
	c.mx.Lock()
	defer c.mx.Unlock()
	if strings.ContainsRune(line, '\n') {
		return nil, ascii.Errorf("at: invalid command: %q", line)
	}

	err := ascii.Fprintf(c.rw, "%s\r\n", line)
	if err != nil {
		return nil, ascii.Errorf("at: write: %w", err)
	}

	var resp Response
	for c.s.Scan() {
		line := c.s.Text()
		switch {
		case line == "OK":
			resp.OK = true
			return &resp, nil
		case line == "ERROR":
			return &resp, nil
		case strings.HasPrefix(line, "+"+name+": "):
			line = strings.TrimPrefix(line, "+"+name+": ")
			resp.Data = append(resp.Data, line)
		default:
			return nil, ascii.Errorf("at: invalid response: %q", line)
		}
	}

	if c.s.Err() != nil {
		return nil, ascii.Errorf("at: read: %w", c.s.Err())
	}

	return nil, io.ErrUnexpectedEOF
}
