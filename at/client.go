package at

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"sync"
)

// Client is an AT command client.
type Client struct {
	rw io.ReadWriter
	s  *bufio.Scanner
	mx sync.Mutex
}

// NewClient creates a new AT command client.
func NewClient(rw io.ReadWriter) *Client {
	s := bufio.NewScanner(rw)
	s.Split(bufio.ScanLines)

	return &Client{rw: rw, s: s}
}

var builders = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
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

	b := builders.Get().(*strings.Builder)
	b.Reset()
	b.WriteString("AT+")
	b.WriteString(name)
	b.WriteString("=")
	for i, param := range params {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(EscapeString(param, ','))
	}
	defer builders.Put(b)

	return c._Command(name, b.String())
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

	return c._Command(name, "AT+"+name+"?")
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

	return c._Command(name, "AT+"+name)
}

func (c *Client) _Command(name, line string) (*Response, error) {
	if strings.ContainsRune(line, '\n') {
		return nil, errors.New("at: invalid command: " + line)
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	_, err := io.WriteString(c.rw, line+"\r\n")
	if err != nil {
		return nil, err
	}

	resp := new(Response)
	for c.s.Scan() {
		line := c.s.Text()

		switch {
		case line == "":
			return nil, io.ErrUnexpectedEOF
		case line == "OK":
			resp.OK = true
			return resp, nil
		case line == "ERROR":
			return resp, nil
		case strings.HasPrefix(line, "+"+name+": "):
			line = strings.TrimPrefix(line, "+"+name+": ")
			resp.Data = append(resp.Data, line)
		default:
			return nil, errors.New("at: invalid response: " + line)
		}
	}

	if c.s.Err() != nil {
		return nil, c.s.Err()
	}

	return nil, io.ErrUnexpectedEOF
}
