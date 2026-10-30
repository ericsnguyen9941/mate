package timeline_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/gitx"
	"github.com/nguyenngocanh94/mate/internal/harness/catalog"
	"github.com/nguyenngocanh94/mate/internal/timeline"
)

// TestASettledCrewsTranscriptIsNotReadAgain is CON-140: every console start
// decoded every transcript a project ever had (421 MB for mihuv2's closed
// crews, ~20 s), though a closed crew's facts were all on the table. The
// transcript is made unreadable after the first pass, size and mtime kept,
// so a pass that reads it fails and a pass that leaves it unread does not.
func TestASettledCrewsTranscriptIsNotReadAgain(t *testing.T) {
	f := newFixture(t)
	transcript := filepath.Join(t.TempDir(), "crew-rollout.jsonl")
	raw, err := os.ReadFile(abs(t, codexFixture))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(transcript, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	f.closeCrew(t, "finished", transcript)
	// A fresh Ingester each time is a console start: nothing is cached.
	cold := func(now time.Time) error {
		return timeline.New(f.ws, f.db, timeline.Deps{
			Now:       func() time.Time { return now },
			Harnesses: catalog.Default(),
			TranscriptRoots: transcriptRoots(filepath.Join(f.root, "no-claude-projects"),
				filepath.Join(f.root, "no-codex-sessions")),
			Git: gitx.Git{Runner: gitx.ExecRunner{}},
		}).Ingest(context.Background())
	}
	if err := cold(fixtureNow); err != nil {
		t.Fatalf("first pass: %v", err)
	}
	turns := `SELECT COUNT(*) FROM turn WHERE actor_id = ?`
	recorded := f.count(t, turns, f.crewActor())
	if recorded == 0 {
		t.Fatal("the fixture crew's transcript recorded no turn")
	}
	if err := os.Chmod(transcript, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(transcript, 0o644) })

	if err := cold(fixtureNow); err == nil {
		t.Fatal("a crew closed moments ago was left unread: its commit sightings may still need git")
	}
	later := fixtureNow.Add(2 * time.Hour)
	if err := cold(later); err != nil {
		t.Fatalf("a crew closed two hours ago was read again: %v", err)
	}
	if n := f.count(t, turns, f.crewActor()); n != recorded {
		t.Fatalf("%d turn(s) after the pass that left the transcript unread, want %d", n, recorded)
	}

	if err := f.ws.AppendStatus(fixtureProject, fixtureCrew, "working: a line after the close"); err != nil {
		t.Fatal(err)
	}
	if err := cold(later); err == nil {
		t.Fatal("a crew with a status line nobody dated was left unread: dating it needs the transcript")
	}
}

// TestAFinalizedMateSessionIsNotReadAgain: an archived Mate session is frozen
// on disk; once the timeline holds it as it is, a console start leaves it
// unread (CON-140). mihuv2 had 59 of them, 112 MB.
func TestAFinalizedMateSessionIsNotReadAgain(t *testing.T) {
	f := newFixture(t)
	archived := filepath.Join(f.root, "archived-transcript.jsonl")
	copyFile(t, claudeFixture, archived)
	if err := f.ws.FreezeMateSession(fixtureProject, map[string]string{
		"harness":    "claude",
		"agent":      "mate-shop",
		"session_id": "archived-0000-4000-8000-000000000000",
		"transcript": archived,
		"started_at": mustTime("2026-09-19T10:40:00Z").Format(time.RFC3339),
	}); err != nil {
		t.Fatalf("FreezeMateSession: %v", err)
	}
	f.ingest(t)
	// The archive reads its own frozen copy, not the live file.
	frozen := filepath.Join(f.ws.MateDir(fixtureProject), "sessions", "archived-0000-4000-8000-000000000000.jsonl")
	if err := os.Chmod(frozen, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(frozen, 0o600) })
	err := timeline.New(f.ws, f.db, timeline.Deps{
		Now:       func() time.Time { return fixtureNow },
		Harnesses: catalog.Default(),
		TranscriptRoots: transcriptRoots(filepath.Join(f.root, "no-claude-projects"),
			filepath.Join(f.root, "no-codex-sessions")),
	}).Ingest(context.Background())
	if err != nil {
		t.Fatalf("a finalized Mate session the timeline holds was read again: %v", err)
	}
}

// TestAConsoleStartAsksGitNothingAboutASettledCrew: a crew closed over an
// hour ago whose merge is recorded costs a cold pass no git call; one whose
// pull request has not landed yet is still asked, however old (CON-140).
func TestAConsoleStartAsksGitNothingAboutASettledCrew(t *testing.T) {
	for _, tc := range []struct {
		name   string
		merged bool
	}{{"merged", true}, {"waiting for its pull request", false}} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFixture(t)
			repo := filepath.Join(f.ws.ProjectHome(fixtureProject), fixtureProject)
			short, _ := commitOnBranch(t, repo, fixtureBranch, "[Buy](pages/checkout-express.html)")
			if tc.merged {
				gitRun(t, repo, "merge", "--ff-only", fixtureBranch)
				gitRun(t, repo, "branch", "-D", fixtureBranch)
			}
			rollout := filepath.Join(f.root, "crew-rollout.jsonl")
			writeCommitRollout(t, rollout, repo, fixtureBranch, short)
			f.closeCrew(t, "finished", rollout)
			f.ingest(t)
			if n := f.count(t, `SELECT COUNT(*) FROM event WHERE kind = ?`, timeline.KindMergeDone); (n == 1) != tc.merged {
				t.Fatalf("%d merge.done event(s), want merged=%v", n, tc.merged)
			}

			counter := &countingGit{}
			later := fixtureNow.Add(2 * time.Hour)
			if err := timeline.New(f.ws, f.db, timeline.Deps{
				Now:       func() time.Time { return later },
				Harnesses: catalog.Default(),
				TranscriptRoots: transcriptRoots(filepath.Join(f.root, "no-claude-projects"),
					filepath.Join(f.root, "no-codex-sessions")),
				Git: gitx.Git{Runner: counter},
			}).Ingest(context.Background()); err != nil {
				t.Fatalf("Ingest: %v", err)
			}
			runs := counter.take()
			if tc.merged && len(runs) != 0 {
				t.Fatalf("a cold pass ran %d git command(s) for a crew settled two hours ago: %q", len(runs), runs)
			}
			if !tc.merged && len(runs) == 0 {
				t.Fatal("a cold pass asked git nothing about a crew whose pull request may still land")
			}
		})
	}
}
