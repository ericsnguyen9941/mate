//go:build darwin || dragonfly || freebsd || illumos || linux || netbsd || openbsd

package platform

import (
	"bytes"
	"os"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

// Usage is what the system has counted of this process (CON-263): the CPU
// time it has used, user and system, its peak resident set, and its
// resident set now where reading it is one small file (Linux's
// /proc/self/statm); RSS is 0 elsewhere. getrusage(2) is one system call,
// cheap enough to ask whenever someone asks.
type Usage struct {
	CPU     time.Duration
	PeakRSS uint64
	RSS     uint64
}

// ProcessUsage is this process's Usage.
func ProcessUsage() (Usage, error) {
	var ru syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &ru); err != nil {
		return Usage{}, err
	}
	u := Usage{CPU: time.Duration(ru.Utime.Nano() + ru.Stime.Nano())}
	// ru_maxrss is bytes on macOS and kilobytes on the others.
	peak := uint64(ru.Maxrss)
	if runtime.GOOS != "darwin" {
		peak *= 1024
	}
	u.PeakRSS = peak
	if runtime.GOOS == "linux" {
		u.RSS = statmRSS()
	}
	return u, nil
}

// statmRSS is the resident set /proc/self/statm says, its second field in
// pages; 0 when it cannot be read.
func statmRSS() uint64 {
	data, err := os.ReadFile("/proc/self/statm")
	if err != nil {
		return 0
	}
	fields := bytes.Fields(data)
	if len(fields) < 2 {
		return 0
	}
	pages, err := strconv.ParseUint(string(fields[1]), 10, 64)
	if err != nil {
		return 0
	}
	return pages * uint64(os.Getpagesize())
}
