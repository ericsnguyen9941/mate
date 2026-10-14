package process

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestExecRunnerCapturesOutputAndExitCode(t *testing.T) {
	t.Parallel()
	r := ExecRunner{}
	res, err := r.Run(context.Background(), Spec{Name: "sh", Args: []string{"-c", "echo out; echo err 1>&2; exit 3"}})
	if err != nil {
		t.Fatalf("non-zero exit must be a result, got error %v", err)
	}
	if res.ExitCode != 3 || strings.TrimSpace(string(res.Stdout)) != "out" || strings.TrimSpace(string(res.Stderr)) != "err" {
		t.Fatalf("result = %+v", res)
	}
}

func TestExecRunnerHonoursDirStdinAndEnv(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	res, err := ExecRunner{}.Run(context.Background(), Spec{
		Name: "sh", Args: []string{"-c", "pwd; cat; echo $MATE_PROBE"}, Dir: dir, Stdin: []byte("in\n"),
		Env: []string{"PATH=/usr/bin:/bin", "MATE_PROBE=yes"},
	})
	if err != nil {
		t.Fatal(err)
	}
	out := string(res.Stdout)
	if !strings.Contains(out, "in\n") || !strings.Contains(out, "yes") {
		t.Fatalf("stdout = %q", out)
	}
	// macOS temp dirs may resolve through /private; pwd prints the resolved form.
	if !strings.Contains(out, strings.TrimPrefix(dir, "/private")) {
		t.Fatalf("stdout %q does not mention dir %q", out, dir)
	}
}

func TestExecRunnerMissingBinaryIsError(t *testing.T) {
	t.Parallel()
	_, err := ExecRunner{}.Run(context.Background(), Spec{Name: "mate-definitely-not-a-binary"})
	if err == nil {
		t.Fatal("missing binary must be an error")
	}
}

func TestExecRunnerReturnsContextError(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := ExecRunner{}.Run(ctx, Spec{Name: "sleep", Args: []string{"5"}})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want deadline exceeded", err)
	}
	cancelled, cancelNow := context.WithCancel(context.Background())
	cancelNow()
	_, err = ExecRunner{}.Run(cancelled, Spec{Name: "true"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want canceled", err)
	}
}

// A grandchild that inherits stdout keeps the pipe open after the process
// itself is killed. The deadline still has to end Run: a caller that bounded
// a Herdr call must not wait for whatever that call left running.
func TestExecRunnerReturnsAtTheDeadlineWhenAGrandchildHoldsStdout(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := ExecRunner{}.Run(ctx, Spec{Name: "sh", Args: []string{"-c", "sleep 5 & wait"}})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want deadline exceeded", err)
	}
	if took := time.Since(start); took > 2*time.Second {
		t.Fatalf("Run returned after %s, want it bounded by the deadline plus WaitDelay", took)
	}
}

// A process that exits 0 and leaves a background child on its stdout has
// answered: Run returns what it wrote instead of waiting for the child.
func TestExecRunnerDoesNotWaitForABackgroundChildAfterASuccessfulExit(t *testing.T) {
	t.Parallel()
	start := time.Now()
	res, err := ExecRunner{}.Run(context.Background(), Spec{Name: "sh", Args: []string{"-c", "echo out; sleep 5 &"}})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if res.ExitCode != 0 || strings.TrimSpace(string(res.Stdout)) != "out" {
		t.Fatalf("result = %+v", res)
	}
	if took := time.Since(start); took > 2*time.Second {
		t.Fatalf("Run returned after %s, want it bounded by WaitDelay", took)
	}
}
