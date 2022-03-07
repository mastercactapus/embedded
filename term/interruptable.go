package term

import "time"

type Interruptable struct {
	fn func(abort func() bool) bool

	lastRun chan bool

	abort   bool
	abortCh chan struct{}
}

func NewInterruptable(runFunc func(abort func() bool) bool) *Interruptable {
	i := &Interruptable{
		fn:      runFunc,
		abortCh: make(chan struct{}, 1),
		lastRun: make(chan bool, 1),
	}
	i.lastRun <- false
	return i
}

func (i *Interruptable) Run() {
	i.Interrupt()
	i.abort = false
	<-i.abortCh
	<-i.lastRun

	go i._run()
}

func (i *Interruptable) shouldAbort() bool {
	if i.abort {
		return true
	}

	select {
	case <-i.abortCh:
		i.abort = true
		return true
	default:
		return false
	}
}

func (i *Interruptable) _run() {
	t := time.NewTimer(100 * time.Millisecond)
	defer t.Stop()
	select {
	case <-i.abortCh:
		// abort during debounce
		i.lastRun <- true
		return
	case <-t.C:
	}
	i.lastRun <- i.fn(i.shouldAbort)
}

func (i *Interruptable) RunSync() {
	i.Interrupt()
	i.abort = false
	<-i.abortCh
	<-i.lastRun

	i._run()
}

// Interrupt returns true if the runFunc was interrupted.
func (i *Interruptable) Interrupt() bool {
	select {
	case i.abortCh <- struct{}{}:
	default:
	}

	last := <-i.lastRun
	i.lastRun <- last
	return last
}
