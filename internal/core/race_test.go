//go:build race

package core

// raceShrink divides the budget test's closed crews under the race
// detector, which runs the ingest some ten times slower: the allocation
// budgets shrink with them (budgetShare), and each call takes about as long
// as without it. raceTime widens the time budgets for what is left of the
// slowdown on a laptop running every package's tests at once.
const (
	raceShrink = 10
	raceTime   = 3
)
