package coreclient

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/coreapi"
)

// A notice (CON-260) reaches the core as POST coreapi.PathChanged with the
// Change as its body.
func TestNotifyPostsTheChange(t *testing.T) {
	got := make(chan coreapi.Change, 1)
	sock := serve(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != coreapi.PathChanged {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		var c coreapi.Change
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			t.Errorf("decode: %v", err)
		}
		got <- c
		w.WriteHeader(http.StatusNoContent)
	})
	want := coreapi.Change{Project: "shop", Crew: "fix-cart", Kind: "status"}
	Notify(sock, want)
	select {
	case c := <-got:
		if c != want {
			t.Fatalf("core got %+v, want %+v", c, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("the core got no notice")
	}
}

// fast fails the test when fn takes longer than bound.
func fast(t *testing.T, what string, bound time.Duration, fn func()) {
	t.Helper()
	start := time.Now()
	fn()
	if took := time.Since(start); took > bound {
		t.Fatalf("%s took %s, want under %s", what, took, bound)
	}
}

// With no core - no socket, or the file a dead core left - a writer pays a
// failed dial and nothing more.
func TestNotifyWithNoCoreIsFast(t *testing.T) {
	dir, err := os.MkdirTemp("/tmp", "cn")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	missing := filepath.Join(dir, "missing.sock")
	fast(t, "100 notices to no socket", 50*time.Millisecond, func() {
		for range 100 {
			Notify(missing, coreapi.Change{Project: "shop", Kind: "sent"})
		}
	})
	// A socket file nothing listens on any more.
	left := filepath.Join(dir, "core.sock")
	ln, err := net.Listen("unix", left)
	if err != nil {
		t.Fatal(err)
	}
	ln.(*net.UnixListener).SetUnlinkOnClose(false)
	_ = ln.Close()
	if _, err := os.Stat(left); err != nil {
		t.Fatalf("the dead core's socket file is gone: %v", err)
	}
	fast(t, "100 notices to a dead core's socket", 50*time.Millisecond, func() {
		for range 100 {
			Notify(left, coreapi.Change{Project: "shop", Kind: "sent"})
		}
	})
}

// A core that never takes a connection (wedged in its accept loop) or
// never answers one costs a writer no more than NotifyWait per notice.
func TestNotifyDoesNotWaitForAWedgedCore(t *testing.T) {
	dir, err := os.MkdirTemp("/tmp", "cn")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	deaf := filepath.Join(dir, "deaf.sock")
	ln, err := net.Listen("unix", deaf)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	fast(t, "20 notices to a core that accepts nothing", 20*NotifyWait, func() {
		for range 20 {
			Notify(deaf, coreapi.Change{Project: "shop", Kind: "sent"})
		}
	})

	hung := make(chan struct{})
	t.Cleanup(func() { close(hung) })
	slow := serve(t, func(http.ResponseWriter, *http.Request) { <-hung })
	fast(t, "20 notices to a core whose handler hangs", 20*NotifyWait, func() {
		for range 20 {
			Notify(slow, coreapi.Change{Project: "shop", Kind: "sent"})
		}
	})
}

// A notice that ran out its deadline makes the process skip that socket
// for a while; one that failed at once (no core) does not.
func TestBackOffOnlyOnATimeout(t *testing.T) {
	sock := fmt.Sprintf("/tmp/cn-backoff-%d.sock", time.Now().UnixNano())
	t.Cleanup(func() { wedged.Delete(sock) })
	backOff(sock, os.ErrNotExist)
	if _, ok := wedged.Load(sock); ok {
		t.Fatal("a failed dial started a backoff")
	}
	backOff(sock, os.ErrDeadlineExceeded)
	if _, ok := wedged.Load(sock); !ok {
		t.Fatal("a deadline that ran out started no backoff")
	}
	got := make(chan struct{}, 1)
	srv := serve(t, func(http.ResponseWriter, *http.Request) { got <- struct{}{} })
	wedged.Store(srv, time.Now().Add(time.Minute))
	t.Cleanup(func() { wedged.Delete(srv) })
	Notify(srv, coreapi.Change{Kind: "workspace"})
	select {
	case <-got:
		t.Fatal("a socket in its backoff was sent a notice")
	case <-time.After(100 * time.Millisecond):
	}
}
