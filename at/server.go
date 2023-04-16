package at

import (
	"bufio"
	"io"
	"strings"

	"github.com/mastercactapus/embedded/term/ascii"
)

// Server is an AT command server.
//
// The server reads AT commands from an io.Reader.
//
// Data sent will be attributed to the last command
// received.
type Server struct {
	rw io.ReadWriter
	s  *bufio.Scanner

	h map[string]HandlerFunc
}

type HandlerFunc func(Cmd) Response

// NewServer creates a new AT command server.
func NewServer(rw io.ReadWriter) *Server {
	s := bufio.NewScanner(rw)

	s.Split(bufio.ScanLines)
	return &Server{rw: rw, s: s}
}

// HandleFunc registers a handler for a command.
//
// The command name is case insensitive, including the "AT+"
func (s *Server) HandleFunc(name string, h HandlerFunc) {
	if s.h == nil {
		s.h = make(map[string]HandlerFunc)
	}
	name = strings.ToUpper(name)
	if _, ok := s.h[strings.ToUpper(name)]; ok {
		panic("at: duplicate handler for " + name)
	}
	s.h[name] = h
}

// Serve serves AT commands.
func (s *Server) Serve() error {
	var c Cmd
	for s.s.Scan() {
		var isSet bool
		var params string
		c.FullName, params, isSet = strings.Cut(s.s.Text(), "=")
		c.FullName = strings.ToUpper(c.FullName)
		c.FullName = strings.TrimSpace(c.FullName)
		if isSet {
			c.FullName += "="
			c.Params = strings.Split(params, ",")
			for i := range c.Params {
				c.Params[i] = UnescapeString(c.Params[i], ',')
			}
		}

		var resp Response
		if c.FullName == "" || c.Name() == "" {
			resp.OK = true
		} else if c.FullName == "START" {
			_, err := s.rw.Write([]byte("HELLO\r\n"))
			if err != nil {
				return err
			}
			continue
		} else if h, ok := s.h[c.FullName]; ok {
			resp = h(c)
		} else {
			resp.SetValue("ERROR", "unknown command")
		}

		if err := s.respond(c, resp); err != nil {
			return err
		}
	}
	if s.s.Err() != nil {
		return s.s.Err()
	}

	return io.ErrUnexpectedEOF
}

func (s *Server) respond(c Cmd, resp Response) error {
	for _, data := range resp.Data {
		if strings.ContainsRune(data, '\n') {
			panic("at: data cannot contain newlines")
		}
		if err := ascii.Fprintf(s.rw, "+%s: %s\r\n", c.Name(), data); err != nil {
			return err
		}
	}

	if resp.OK {
		if err := ascii.Fprintf(s.rw, "OK\r\n"); err != nil {
			return err
		}
	} else {
		if err := ascii.Fprintf(s.rw, "ERROR\r\n"); err != nil {
			return err
		}
	}

	return nil
}
