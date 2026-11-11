package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"time"
)

// pprofEnv names the address `mate core run` serves Go's profiles on
// (CON-257): `MATE_PPROF=127.0.0.1:6060`. Unset, nothing is served. The
// investigation of 2026-11-11 could only profile a copy of the workspace
// beside the running core; with it, `go tool pprof
// http://127.0.0.1:6060/debug/pprof/heap` reads the core itself.
const pprofEnv = "MATE_PPROF"

// startPprof serves /debug/pprof/ on addr until stop is called. Only a
// loopback address is taken: the profiles show the workspace's paths and
// whatever strings the heap holds, and nothing guards the port.
func startPprof(addr string) (bound string, stop func(), err error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", nil, fmt.Errorf("%s=%q: %w", pprofEnv, addr, err)
	}
	if ip := net.ParseIP(host); host != "localhost" && (ip == nil || !ip.IsLoopback()) {
		return "", nil, fmt.Errorf("%s=%q: only a loopback address is served", pprofEnv, addr)
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return "", nil, fmt.Errorf("%s: %w", pprofEnv, err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() { _ = srv.Serve(ln) }()
	return ln.Addr().String(), func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}, nil
}
