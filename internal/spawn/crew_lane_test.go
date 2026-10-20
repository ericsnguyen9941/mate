package spawn

import (
	"testing"

	"github.com/nguyenngocanh94/mate/internal/store"
)

// A spawn holds its repo's lane unless --lane says none or names one
// (CON-118); a lane is one word.
func TestSpawnLaneDefaultsToTheRepos(t *testing.T) {
	repo := store.RepoConfig{Name: "app", Lane: "xcode"}
	for _, tc := range []struct{ flag, want string }{{"", "xcode"}, {"none", ""}, {"gpu", "gpu"}} {
		got, err := spawnLane(SpawnCrewRequest{Lane: tc.flag}, repo)
		if err != nil || got != tc.want {
			t.Errorf("--lane %q = %q, %v; want %q", tc.flag, got, err, tc.want)
		}
	}
	if got, err := spawnLane(SpawnCrewRequest{}, store.RepoConfig{Name: "api"}); err != nil || got != "" {
		t.Errorf("a repo with no lane gave %q, %v", got, err)
	}
	if _, err := spawnLane(SpawnCrewRequest{Lane: "two words"}, repo); err == nil {
		t.Error("a two-word lane was accepted")
	}
}
