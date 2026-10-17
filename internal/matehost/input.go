package matehost

import (
	"errors"
	"os"
	"time"

	"github.com/charmbracelet/x/term"
	"golang.org/x/sys/unix"
)

// Bubble Tea v1 reads input in chunks and parses each chunk on its own: a
// chunk that ends inside an escape sequence has its head read as Alt+[ and
// its tail typed into the composer. A Mate pane gets its input through
// Herdr's pty, where a burst of wheel notches is split anywhere, so
// `<64;35;39M` landed in the composer (2026-10-17). seqInput hands Bubble
// Tea only chunks that end on a sequence boundary.

// seqWait is how long a chunk that ends inside an escape sequence waits for
// the rest of it. A lone Esc key waits this long too.
const seqWait = 50 * time.Millisecond

// seqRoom is what a read leaves free in the caller's buffer for the rest of
// a sequence: an SGR mouse report is at most 16 bytes.
const seqRoom = 32

// seqInput is a terminal's input that never returns a chunk ending inside
// an escape sequence. It keeps the file's descriptor and name, so Bubble
// Tea still puts the terminal in raw mode and can cancel a read.
type seqInput struct {
	*os.File
}

// newSeqInput wraps f, or returns nil when f is not a terminal: Bubble
// Tea's own input handling is right for a pipe.
func newSeqInput(f *os.File) *seqInput {
	if !term.IsTerminal(f.Fd()) {
		return nil
	}
	return &seqInput{f}
}

func (s *seqInput) Read(p []byte) (int, error) {
	limit := len(p)
	if limit > seqRoom*2 {
		limit -= seqRoom
	}
	n, err := s.File.Read(p[:limit])
	for err == nil && n < len(p) && openSequence(p[:n]) && s.readable(seqWait) {
		var m int
		m, err = s.File.Read(p[n:])
		n += m
	}
	return n, err
}

// readable waits up to d for the file to have input.
func (s *seqInput) readable(d time.Duration) bool {
	fds := []unix.PollFd{{Fd: int32(s.Fd()), Events: unix.POLLIN}}
	for {
		n, err := unix.Poll(fds, int(d/time.Millisecond))
		if errors.Is(err, unix.EINTR) {
			continue
		}
		return err == nil && n > 0
	}
}

// openSequence reports whether b ends inside an escape sequence: a lone
// Esc, Esc O, or a CSI without its final byte.
func openSequence(b []byte) bool {
	i := len(b) - 1
	for i >= 0 && b[i] != 0x1b {
		i--
	}
	if i < 0 {
		return false
	}
	tail := b[i+1:]
	if len(tail) == 0 {
		return true
	}
	switch tail[0] {
	case 'O':
		return len(tail) == 1
	case '[':
		for _, c := range tail[1:] {
			if c < 0x20 || c > 0x3f {
				// A final byte (0x40-0x7e) ends it; anything else is
				// not a CSI and is Bubble Tea's to judge.
				return false
			}
		}
		return true
	}
	return false
}
