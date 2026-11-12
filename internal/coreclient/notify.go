package coreclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
)

// NotifyWait bounds one notice: the dial and the write together. A core
// that runs accepts at once and the request fits the socket's buffer, so a
// notice costs a writer well under a millisecond; one with no core costs a
// failed dial, as fast.
const NotifyWait = 100 * time.Millisecond

// notifyBackoff is how long a socket that ran out NotifyWait gets no
// notice from this process: a wedged core costs a writer one NotifyWait,
// not one per write.
const notifyBackoff = 10 * time.Second

// wedged is each socket's end of backoff, in this process.
var wedged sync.Map // socket path -> time.Time

// Notify tells the core on socket that a record changed (POST
// coreapi.PathChanged, CON-260), and returns without the answer: the
// request is written and the connection closed. It never fails and never
// waits more than NotifyWait: with no core (no socket, or a socket a dead
// core left) the dial fails at once, and a core that does not take the
// request in time is skipped by this process for the next 10 seconds. A
// notice is a hint; the core's polls find what one did not say.
//
// cmd/mate makes it the store's ChangeNotifier in every process but the
// core's own.
func Notify(socket string, c coreapi.Change) {
	if until, ok := wedged.Load(socket); ok && time.Now().Before(until.(time.Time)) {
		return
	}
	body, err := json.Marshal(c)
	if err != nil {
		return
	}
	deadline := time.Now().Add(NotifyWait)
	conn, err := net.DialTimeout("unix", socket, NotifyWait)
	if err != nil {
		backOff(socket, err)
		return
	}
	defer conn.Close()
	_ = conn.SetDeadline(deadline)
	req, err := http.NewRequest(http.MethodPost, "http://core"+coreapi.PathChanged, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	var buf bytes.Buffer
	if err := req.Write(&buf); err != nil {
		return
	}
	// One write: the whole request lands in the socket's buffer or the
	// deadline ends it.
	_, err = conn.Write(buf.Bytes())
	backOff(socket, err)
}

// backOff starts the socket's backoff when err is a deadline that ran out;
// any other failure (no core) costs nothing to try again.
func backOff(socket string, err error) {
	if err == nil {
		return
	}
	var ne net.Error
	if errors.Is(err, os.ErrDeadlineExceeded) || (errors.As(err, &ne) && ne.Timeout()) {
		wedged.Store(socket, time.Now().Add(notifyBackoff))
	}
}
