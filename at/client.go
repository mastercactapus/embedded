package at

import (
	"bufio"
	"context"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/mastercactapus/embedded/term/ascii"
)

// Client is an AT command client.
type Client struct {
	rw      io.ReadWriter
	sr      *ScanReader
	mx      sync.Mutex
	timeout time.Duration
}

// NewClient creates a new AT command client.
func NewClient(rw io.ReadWriter) *Client {
	s := bufio.NewScanner(rw)
	s.Split(bufio.ScanLines)

	return &Client{rw: rw, sr: NewScanReader(s)}
}

// SetTimeout sets the timeout for reading command responses
// from the modem.
func (c *Client) SetTimeout(d time.Duration) {
	c.timeout = d
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

func (c *Client) _Command(name, line string) (*Response, error) {
	c.mx.Lock()
	defer c.mx.Unlock()
	if strings.ContainsRune(line, '\n') {
		return nil, ascii.Errorf("at: invalid command: %q", line)
	}

	_, err := io.WriteString(c.rw, line+"\r\n")
	if err != nil {
		return nil, ascii.Errorf("at: write: %w", err)
	}
	ctx := context.Background()
	if c.timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	var resp Response
	for {
		line, err := c.sr.Next(ctx)
		if err != nil {
			return nil, ascii.Errorf("at: read: %w", err)
		}

		switch {
		case line == "":
			return nil, io.ErrUnexpectedEOF
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
}
