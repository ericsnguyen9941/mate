//go:build darwin || dragonfly || freebsd || illumos || linux || netbsd || openbsd

package platform

import (
	"runtime"
	"testing"
)

func TestProcessUsageCountsThisProcess(t *testing.T) {
	// Some CPU to count.
	sum := 0
	for i := range 5_000_000 {
		sum += i
	}
	_ = sum
	u, err := ProcessUsage()
	if err != nil {
		t.Fatal(err)
	}
	if u.CPU <= 0 {
		t.Errorf("CPU = %s, want the time this test used", u.CPU)
	}
	// A Go test binary's resident set is megabytes, not kilobytes: a
	// kilobyte count left unscaled would read as under 1 MB.
	if u.PeakRSS < 1<<20 {
		t.Errorf("PeakRSS = %d bytes, want at least 1 MB", u.PeakRSS)
	}
	if runtime.GOOS == "linux" && u.RSS < 1<<20 {
		t.Errorf("RSS = %d bytes, want at least 1 MB on Linux", u.RSS)
	}
}
