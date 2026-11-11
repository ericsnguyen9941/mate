//go:build !darwin && !dragonfly && !freebsd && !illumos && !linux && !netbsd && !openbsd

package platform

import "time"

// Usage is what the system has counted of this process; see usage_unix.go.
type Usage struct {
	CPU     time.Duration
	PeakRSS uint64
	RSS     uint64
}

// errNoUsage is ProcessUsage's answer here: nothing was counted.
const errNoUsage unsupported = "reading a process's resource usage is not supported"

func ProcessUsage() (Usage, error) { return Usage{}, errNoUsage }
