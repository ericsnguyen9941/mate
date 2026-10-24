package matehost

import (
	"testing"
	"time"
)

// The forwarder hands events on in order and never makes the host wait for
// the program: here the program is busy (as it is while Update runs) until
// after the host's call has returned.
func TestForwarderKeepsOrderWithoutBlockingTheHost(t *testing.T) {
	got := make(chan Event, 8)
	busy := make(chan struct{})
	f := newForwarder(func(ev Event) {
		<-busy
		got <- ev
	})
	defer f.Close()
	done := make(chan struct{})
	go func() {
		for i := 1; i <= 3; i++ {
			f.Put(Event{Queue: []string{string(rune('0' + i))}})
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Put waited for the program")
	}
	close(busy)
	for i := 1; i <= 3; i++ {
		select {
		case ev := <-got:
			if want := string(rune('0' + i)); ev.Queue[0] != want {
				t.Fatalf("event %d: got %q, want %q", i, ev.Queue[0], want)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("event %d never arrived", i)
		}
	}
}

// Opening a session from the REPL's Update emits an event; with the program
// busy in that very Update, the call still returns (CON-121).
func TestOpenSessionReturnsWhileTheProgramIsBusy(t *testing.T) {
	f := newSessionFixture(t)
	inUpdate := make(chan struct{})
	seen := make(chan Event, 16)
	fw := newForwarder(func(ev Event) {
		<-inUpdate
		seen <- ev
	})
	defer fw.Close()
	f.host.Subscribe(fw.Put)
	done := make(chan error, 1)
	go func() { done <- f.host.OpenSession(1) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("OpenSession waited for the program's event loop")
	}
	close(inUpdate)
	select {
	case ev := <-seen:
		if len(ev.Sessions) != 1 || ev.Sessions[0].News != 1 {
			t.Fatalf("event sessions = %+v, want session of #1", ev.Sessions)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the session's event never reached the program")
	}
}
