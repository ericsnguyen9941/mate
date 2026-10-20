package teardown

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/process"
)

const wt = "/ws/.worktrees/plantrip-pla-6"

func TestNamedAfterNeedsTheWholeWorktreeName(t *testing.T) {
	tg := Target{Worktree: wt}
	for name, want := range map[string]bool{
		"PlanTrip-plantrip-pla-6":          true,
		"plantrip-pla-6":                   true,
		"HV plantrip-pla-6 iPhone 17":      true,
		"PlanTrip-plantrip-pla-6-be-r2":    false,
		"PlanTrip-plantrip-pla-60":         false,
		"PlanTrip-xplantrip-pla-6":         false,
		"PlanTrip-base":                    false,
		"plantrip-pla-6.old":               false,
		"PlanTrip-plantrip-pla-6 (backup)": true,
	} {
		if got := tg.NamedAfter(name); got != want {
			t.Errorf("NamedAfter(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestOwnsOnlyTheWorktreeAndBelow(t *testing.T) {
	tg := Target{Worktree: wt}
	for path, want := range map[string]bool{
		wt:                             true,
		wt + "/HelloVietnam.xcodeproj": true,
		wt + "-be-r2":                  false,
		"/ws/plantrip/plantrip_api":    false,
		"":                             false,
	} {
		if got := tg.Owns(path); got != want {
			t.Errorf("Owns(%q) = %v, want %v", path, got, want)
		}
	}
}

// docker answers like the real CLI for a machine with the crew's stack, a
// reviewer's left-over volume set and the captain's own checkout stack.
func dockerFake(t *testing.T, removed *[]string) *process.FakeRunner {
	return &process.FakeRunner{Handler: func(_ context.Context, s process.Spec) (process.Result, error) {
		if s.Name != "docker" {
			t.Fatalf("unexpected command %s", process.CommandKey(s))
		}
		args := strings.Join(s.Args, " ")
		switch {
		case strings.HasPrefix(args, "ps -a --filter label=com.docker.compose.project --format"):
			return process.Result{Stdout: []byte(
				"plantrip-pla-6\t" + wt + "\n" +
					"plantrip-pla-6\t" + wt + "\n" +
					"plantrip_api\t/ws/plantrip/plantrip_api\n")}, nil
		case strings.HasPrefix(args, "ps -aq --filter label=com.docker.compose.project=plantrip-pla-6"):
			return process.Result{Stdout: []byte("c1\nc2\n")}, nil
		case strings.HasPrefix(args, "volume ls -q --filter label=com.docker.compose.project=plantrip-pla-6"):
			return process.Result{Stdout: []byte("plantrip-pla-6_sail-pgsql\n")}, nil
		case strings.HasPrefix(args, "network ls -q --filter label=com.docker.compose.project=plantrip-pla-6"):
			return process.Result{Stdout: []byte("n1\n")}, nil
		case strings.HasPrefix(args, "rm -f"), strings.HasPrefix(args, "volume rm"), strings.HasPrefix(args, "network rm"):
			*removed = append(*removed, args)
			return process.Result{}, nil
		case strings.Contains(args, "label=com.docker.compose.project="):
			return process.Result{}, nil
		}
		t.Fatalf("unexpected docker %s", args)
		return process.Result{}, nil
	}}
}

func TestDockerRemovesTheWorktreesStackByLabelAndLeavesOthers(t *testing.T) {
	var removed []string
	got, err := Docker{Runner: dockerFake(t, &removed)}.Sweep(context.Background(), Target{Worktree: wt})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"docker plantrip-pla-6 (2 containers, 1 volume, 1 network)"}; !slices.Equal(got, want) {
		t.Fatalf("removed = %q, want %q", got, want)
	}
	if want := []string{"rm -f c1 c2", "volume rm -f plantrip-pla-6_sail-pgsql", "network rm n1"}; !slices.Equal(removed, want) {
		t.Fatalf("docker calls = %q, want %q", removed, want)
	}
}

func TestDockerFindsAReviewersVolumesByTheWorktreeName(t *testing.T) {
	// The reviewer's containers are gone; only its volumes carry the
	// compose project compose derived from the worktree's directory.
	run := &process.FakeRunner{Handler: func(_ context.Context, s process.Spec) (process.Result, error) {
		args := strings.Join(s.Args, " ")
		switch {
		case strings.HasPrefix(args, "ps -a --filter"):
			return process.Result{}, nil
		case strings.HasPrefix(args, "volume ls -q --filter label=com.docker.compose.project=plantrip-pla-31-r1b"):
			return process.Result{Stdout: []byte("plantrip-pla-31-r1b_sail-pgsql\nplantrip-pla-31-r1b_sail-redis\n")}, nil
		}
		return process.Result{}, nil
	}}
	got, err := Docker{Runner: run}.Sweep(context.Background(), Target{Worktree: "/ws/.worktrees/plantrip-pla-31-r1b"})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"docker plantrip-pla-31-r1b (2 volumes)"}; !slices.Equal(got, want) {
		t.Fatalf("removed = %q, want %q", got, want)
	}
}

func TestDockerLeavesAProjectWithAContainerElsewhere(t *testing.T) {
	run := &process.FakeRunner{Handler: func(_ context.Context, s process.Spec) (process.Result, error) {
		args := strings.Join(s.Args, " ")
		if strings.HasPrefix(args, "ps -a --filter label=com.docker.compose.project --format") {
			// Someone else's stack happens to carry the worktree's name.
			return process.Result{Stdout: []byte("plantrip-pla-6\t/elsewhere/plantrip-pla-6\n")}, nil
		}
		t.Fatalf("touched a foreign project: docker %s", args)
		return process.Result{}, nil
	}}
	got, err := Docker{Runner: run}.Sweep(context.Background(), Target{Worktree: wt})
	if err != nil || len(got) != 0 {
		t.Fatalf("removed %q, err %v; want nothing", got, err)
	}
}

func TestDockerNotInstalledIsNothingAndDaemonDownIsPending(t *testing.T) {
	notFound := &process.FakeRunner{Handler: func(context.Context, process.Spec) (process.Result, error) {
		return process.Result{}, fmt.Errorf("process docker: %w", exec.ErrNotFound)
	}}
	if got, err := (Docker{Runner: notFound}).Sweep(context.Background(), Target{Worktree: wt}); err != nil || got != nil {
		t.Fatalf("no docker: %q, %v; want nothing", got, err)
	}
	down := &process.FakeRunner{Default: process.Result{ExitCode: 1, Stderr: []byte("Cannot connect to the Docker daemon")}}
	_, err := Docker{Runner: down}.Sweep(context.Background(), Target{Worktree: wt})
	if err == nil || !strings.Contains(err.Error(), "Cannot connect to the Docker daemon") {
		t.Fatalf("daemon down: err %v, want docker's own words", err)
	}
}

const simJSON = `{"devices":{"com.apple.CoreSimulator.SimRuntime.iOS-27-0":[
 {"name":"PlanTrip-base","udid":"BASE","state":"Shutdown"},
 {"name":"PlanTrip-plantrip-pla-6","udid":"MINE","state":"Booted"},
 {"name":"PlanTrip-plantrip-pla-6-be-r2","udid":"OTHER","state":"Shutdown"},
 {"name":"iPhone 17","udid":"STOCK","state":"Shutdown"}]}}`

func TestSimulatorsDeletesOnlyTheOneNamedAfterTheWorktree(t *testing.T) {
	run := &process.FakeRunner{Responses: map[string]process.Result{
		"xcrun simctl list devices -j": {Stdout: []byte(simJSON)},
	}}
	got, err := Simulators{Runner: run}.Sweep(context.Background(), Target{Worktree: wt})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"simulator PlanTrip-plantrip-pla-6"}; !slices.Equal(got, want) {
		t.Fatalf("removed = %q, want %q", got, want)
	}
	var calls []string
	for _, c := range run.Calls {
		calls = append(calls, process.CommandKey(c))
	}
	if want := []string{"xcrun simctl list devices -j", "xcrun simctl shutdown MINE", "xcrun simctl delete MINE"}; !slices.Equal(calls, want) {
		t.Fatalf("calls = %q, want %q", calls, want)
	}
}

func TestDerivedDataDeletesTheFoldersXcodeMadeForTheWorktree(t *testing.T) {
	root := t.TempDir()
	write := func(name, workspace string) {
		dir := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Join(dir, "Build"), 0o755); err != nil {
			t.Fatal(err)
		}
		plist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>LastAccessedDate</key>
	<date>2026-10-15T02:17:11Z</date>
	<key>WorkspacePath</key>
	<string>` + workspace + `</string>
</dict>
</plist>
`
		if err := os.WriteFile(filepath.Join(dir, "info.plist"), []byte(plist), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("HelloVietnam-mine", wt+"/HelloVietnam.xcodeproj")
	write("HelloVietnam-other", wt+"-be-r2/HelloVietnam.xcodeproj")
	write("HelloVietnam-primary", "/ws/hellovietnam/HelloVietnam.xcodeproj")
	if err := os.MkdirAll(filepath.Join(root, "CompilationCache.noindex"), 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := DerivedData{Root: root}.Sweep(context.Background(), Target{Worktree: wt})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"DerivedData HelloVietnam-mine"}; !slices.Equal(got, want) {
		t.Fatalf("removed = %q, want %q", got, want)
	}
	for name, want := range map[string]bool{"HelloVietnam-mine": false, "HelloVietnam-other": true, "HelloVietnam-primary": true, "CompilationCache.noindex": true} {
		if _, err := os.Stat(filepath.Join(root, name)); (err == nil) != want {
			t.Errorf("%s exists = %v, want %v", name, err == nil, want)
		}
	}
}

func TestProcessesStopsWhatRunsInTheWorktreeButNotAShell(t *testing.T) {
	lsof := "p100\ncphp\nfcwd\nn" + wt + "\n" +
		"p101\nczsh\nfcwd\nn" + wt + "/app\n" +
		"p102\ncnode\nfcwd\nn" + wt + "-be-r2\n" +
		"p103\ncvite\nfcwd\nn" + wt + "/web\n"
	run := &process.FakeRunner{Default: process.Result{Stdout: []byte(lsof)}}
	alive := map[int]bool{100: true, 103: true}
	var sent []string
	signal := func(pid int, sig syscall.Signal) error {
		if !alive[pid] {
			return syscall.ESRCH
		}
		if sig != 0 {
			sent = append(sent, fmt.Sprintf("%d:%d", pid, sig))
		}
		if sig == syscall.SIGTERM && pid == 100 {
			alive[pid] = false // php exits on TERM
		}
		if sig == syscall.SIGKILL {
			alive[pid] = false
		}
		return nil
	}
	got, err := Processes{Runner: run, Signal: signal, Grace: 50 * time.Millisecond}.Sweep(context.Background(), Target{Worktree: wt})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"process php 100", "process vite 103"}; !slices.Equal(got, want) {
		t.Fatalf("removed = %q, want %q", got, want)
	}
	want := []string{fmt.Sprintf("100:%d", syscall.SIGTERM), fmt.Sprintf("103:%d", syscall.SIGTERM), fmt.Sprintf("103:%d", syscall.SIGKILL)}
	if !slices.Equal(sent, want) {
		t.Fatalf("signals = %q, want %q", sent, want)
	}
}

type sweepFunc func(context.Context, Target) ([]string, error)

func (f sweepFunc) Sweep(ctx context.Context, t Target) ([]string, error) { return f(ctx, t) }

func TestCleanRunsTheHookInTheWorktreeThenEverySweep(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "crew", "teardown.log")
	var order []string
	run := &process.FakeRunner{Handler: func(_ context.Context, s process.Spec) (process.Result, error) {
		order = append(order, "hook")
		if s.Name != "sh" || !slices.Equal(s.Args, []string{"-c", "make down"}) || s.Dir != dir {
			t.Fatalf("hook ran as %s in %s", process.CommandKey(s), s.Dir)
		}
		if !slices.Contains(s.Env, "MATE_CREW=pla-6") {
			t.Fatalf("hook env lacks MATE_CREW: %q", s.Env)
		}
		return process.Result{Stdout: []byte("stack down\n")}, nil
	}}
	sweep := func(name string, err error) Sweep {
		return sweepFunc(func(context.Context, Target) ([]string, error) {
			order = append(order, name)
			return []string{name + " gone"}, err
		})
	}
	c := Cleaner{Runner: run, Sweeps: []Sweep{sweep("a", nil), sweep("b", errors.New("b: daemon down")), sweep("c", nil)}}
	r := c.Clean(context.Background(), Target{Worktree: dir, Hook: "make down", Env: []string{"MATE_CREW=pla-6"}, Log: logPath})

	if want := []string{"hook", "a", "b", "c"}; !slices.Equal(order, want) {
		t.Fatalf("order = %q, want %q", order, want)
	}
	if r.Hook != "ok" || !slices.Equal(r.Pending, []string{"b: daemon down"}) {
		t.Fatalf("report = %+v", r)
	}
	if got, want := r.Meta(), "pending: b: daemon down (hook ok, a gone, b gone, c gone)"; got != want {
		t.Fatalf("Meta() = %q, want %q", got, want)
	}
	log, err := os.ReadFile(logPath)
	if err != nil || !strings.Contains(string(log), "$ make down\nstack down") {
		t.Fatalf("log = %q, %v", log, err)
	}
}

func TestCleanReportsAFailedHookAndStillSweeps(t *testing.T) {
	run := &process.FakeRunner{Default: process.Result{ExitCode: 2, Stderr: []byte("make: *** No rule to make target 'down'.")}}
	swept := false
	c := Cleaner{Runner: run, Sweeps: []Sweep{sweepFunc(func(context.Context, Target) ([]string, error) { swept = true; return nil, nil })}}
	r := c.Clean(context.Background(), Target{Worktree: t.TempDir(), Hook: "make down"})
	if !swept || r.Hook != "failed" || len(r.Pending) != 1 || !strings.Contains(r.Pending[0], "No rule to make target") {
		t.Fatalf("swept %v, report %+v", swept, r)
	}
}

func TestCleanSkipsTheHookWhenTheWorktreeIsGone(t *testing.T) {
	run := &process.FakeRunner{}
	r := Cleaner{Runner: run}.Clean(context.Background(), Target{Worktree: filepath.Join(t.TempDir(), "gone"), Hook: "make down"})
	if len(run.Calls) != 0 || r.Ran() {
		t.Fatalf("calls %d, report %+v; want nothing", len(run.Calls), r)
	}
	if got := r.Meta(); got != "done" {
		t.Fatalf("Meta() = %q, want done", got)
	}
}

func TestCleanStopsAHookAtItsTimeout(t *testing.T) {
	c := Cleaner{Runner: process.ExecRunner{}, HookTimeout: 200 * time.Millisecond}
	start := time.Now()
	r := c.Clean(context.Background(), Target{Worktree: t.TempDir(), Hook: "sleep 30"})
	if time.Since(start) > 10*time.Second || len(r.Pending) != 1 || !strings.Contains(r.Pending[0], "did not finish") {
		t.Fatalf("report %+v after %s", r, time.Since(start))
	}
}

func TestComposeProjectNameFollowsCompose(t *testing.T) {
	if got := ComposeProjectName("PlanTrip.pla-6_x"); got != "plantrippla-6_x" {
		t.Fatalf("got %q", got)
	}
}
