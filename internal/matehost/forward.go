package matehost

import "sync"

// forwarder hands the host's events to the program in the order they came,
// without the host ever waiting on the program. tea.Program.Send blocks
// until the event loop takes the message, and the event loop is the
// goroutine that runs Update, so a host call made from Update that emits
// (OpenSession, AskSession, CloseSession) would wait on itself for ever
// (CON-121). The queue is unbounded: the host emits a few events per turn,
// and dropping or reordering them would leave the REPL showing stale state.
type forwarder struct {
	mu    sync.Mutex
	cond  *sync.Cond
	queue []Event
	send  func(Event)
}

// newForwarder starts the goroutine that drains the queue into send; it
// ends when Close is called.
func newForwarder(send func(Event)) *forwarder {
	f := &forwarder{send: send}
	f.cond = sync.NewCond(&f.mu)
	go f.run()
	return f
}

// Put queues ev and returns at once.
func (f *forwarder) Put(ev Event) {
	f.mu.Lock()
	f.queue = append(f.queue, ev)
	f.cond.Signal()
	f.mu.Unlock()
}

// Close stops the goroutine; what is still queued is dropped.
func (f *forwarder) Close() {
	f.mu.Lock()
	f.send = nil
	f.cond.Signal()
	f.mu.Unlock()
}

func (f *forwarder) run() {
	for {
		f.mu.Lock()
		for len(f.queue) == 0 && f.send != nil {
			f.cond.Wait()
		}
		send := f.send
		if send == nil {
			f.mu.Unlock()
			return
		}
		ev := f.queue[0]
		f.queue = f.queue[1:]
		f.mu.Unlock()
		send(ev)
	}
}
