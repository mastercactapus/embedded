package term

import "time"

type Ticker struct {
	C    <-chan time.Time
	stop chan struct{}
}

func NewTicker(dur time.Duration) *Ticker {
	if dur <= 0 {
		panic("invalid duration")
	}

	c := make(chan time.Time, 1)
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-stop:
				return
			case c <- time.Now():
			}
			time.Sleep(dur)
		}
	}()

	return &Ticker{C: c, stop: stop}
}

func (t *Ticker) Stop() { close(t.stop) }
