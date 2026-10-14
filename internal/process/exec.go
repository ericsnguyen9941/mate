package process

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

// WaitDelay is how long Run waits for a process's output pipes once the
// process itself is gone - killed at the context's deadline, or exited. A
// grandchild that inherited stdout would otherwise keep Run waiting after a
// caller's deadline for as long as that grandchild lives, which makes the
// deadline a suggestion.
const WaitDelay = time.Second

// ExecRunner runs processes with os/exec. A non-zero exit is a result, not
// an error: callers read ExitCode. Failing to start the process (missing
// binary, bad directory) is an error. A cancelled or expired context returns
// ctx.Err(), per the Runner contract. A process that exits 0 and leaves a
// background child holding its output has answered: Run returns what it
// wrote within WaitDelay.
type ExecRunner struct{}

// Run implements Runner.
func (ExecRunner) Run(ctx context.Context, spec Spec) (Result, error) {
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if spec.Name == "" {
		return Result{}, errors.New("process: name is required")
	}
	cmd := exec.CommandContext(ctx, spec.Name, spec.Args...)
	cmd.Dir = spec.Dir
	if spec.Env != nil {
		cmd.Env = spec.Env
	}
	if spec.Stdin != nil {
		cmd.Stdin = bytes.NewReader(spec.Stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.WaitDelay = WaitDelay
	err := cmd.Run()
	if errors.Is(err, exec.ErrWaitDelay) {
		// The process exited 0; only a child it left behind still held
		// the pipe.
		err = nil
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		return Result{}, ctxErr
	}
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return Result{ExitCode: exitErr.ExitCode(), Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}, nil
		}
		return Result{}, fmt.Errorf("process %s: %w", spec.Name, err)
	}
	return Result{ExitCode: 0, Stdout: stdout.Bytes(), Stderr: stderr.Bytes()}, nil
}
