package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

// MATE_PPROF serves Go's profiles on a loopback address and refuses any
// other: nothing guards the port (CON-257).
func TestStartPprofServesOnLoopbackOnly(t *testing.T) {
	for _, addr := range []string{"0.0.0.0:0", ":0", "192.168.1.5:6060", "example.com:6060", "6060"} {
		if _, stop, err := startPprof(addr); err == nil {
			stop()
			t.Fatalf("startPprof(%q) served; want a refusal", addr)
		}
	}
	bound, stop, err := startPprof("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer stop()
	resp, err := http.Get("http://" + bound + "/debug/pprof/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "heap") {
		t.Fatalf("index: %d %s", resp.StatusCode, body)
	}
}
