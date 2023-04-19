package at

import (
	"bufio"
	"io"
	"strings"
	"sync"
	"time"
)

// Server is an AT command server.
//
// The server reads AT commands from an io.Reader.
//
// Data sent will be attributed to the last command
// received.
type Server struct {
	s *bufio.Scanner
	w *bufio.Writer

	idleDur time.Duration
	t       *time.Timer
	mx      sync.Mutex

	h map[string]HandlerFunc
}

type HandlerFunc func(Cmd) Response

// NewServer creates a new AT command server.
func NewServer(rw io.ReadWriter) *Server {
	s := bufio.NewScanner(rw)

	s.Split(bufio.ScanLines)
	return &Server{w: bufio.NewWriter(rw), s: s}
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

// OnIdle sets a handler to be called when the server
// has been idle for the specified duration.
//
// OnIdle will panic if called more than once.
//
// It is guaranteed that the handler will not be called
// while a command is being handled.
func (s *Server) OnIdle(dur time.Duration, h func()) {
	if s.t != nil {
		panic("at: idle handler already set")
	}
	s.idleDur = dur
	s.t = time.AfterFunc(dur, func() {
		s.mx.Lock()
		h()
		s.mx.Unlock()
	})
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

		if s.t != nil {
			s.t.Reset(s.idleDur)
		}

		var resp Response
		if c.FullName == "" || c.Name() == "" {
			resp.OK = true
		} else if h, ok := s.h[c.FullName]; ok {
			s.mx.Lock()
			resp = h(c)
			s.mx.Unlock()
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

		if _, err := io.WriteString(s.w, "+"+c.Name()+": "+data+"\r\n"); err != nil {
			return err
		}
	}

	if resp.OK {
		if _, err := io.WriteString(s.w, "OK\r\n"); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(s.w, "ERROR\r\n"); err != nil {
			return err
		}
	}

	return s.w.Flush()
}
