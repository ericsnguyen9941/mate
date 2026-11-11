//go:build !race

package core

// Without the race detector the budget test runs as written (race_test.go).
const (
	raceShrink = 1
	raceTime   = 1
)
