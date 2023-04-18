package at

import (
	"context"
)

type Scanner interface {
	Scan() bool
	Text() string
	Err() error
}
type ScanReader struct {
	s Scanner

	scanCh chan bool
	textCh chan string
}

func NewScanReader(s Scanner) *ScanReader {
	r := &ScanReader{
		s:      s,
		scanCh: make(chan bool),
		textCh: make(chan string),
	}

	go func() {
		for {
			res := r.s.Scan()
			r.scanCh <- res
			if !res {
				close(r.textCh)
				close(r.scanCh)
				return
			}
			r.textCh <- r.s.Text()
		}
	}()

	return r
}

func (s *ScanReader) Next(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-s.scanCh:
		if !res {
			return "", s.s.Err()
		}
		return <-s.textCh, nil
	}
}
