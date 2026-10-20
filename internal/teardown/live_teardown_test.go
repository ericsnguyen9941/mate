package teardown

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestLiveTeardownRemovesWhatAWorktreeStarted is the proof on a real
// machine (CON-111): a compose-labelled container, volume and network
// started from a worktree, a simulator named after it, an Xcode DerivedData
// folder recorded for it and a process running in it are all gone after
// one Clean, and a container from another directory is not.
func TestLiveTeardownRemovesWhatAWorktreeStarted(t *testing.T) {
	if os.Getenv("MATE_LIVE") != "1" {
		t.Skip("set MATE_LIVE=1 to run live teardown proofs (docker, xcrun simctl)")
	}
	ctx := context.Background()
	parent := t.TempDir()
	name := "mate-live-teardown-" + time.Now().Format("150405")
	worktree := filepath.Join(parent, name)
	if err := os.MkdirAll(worktree, 0o755); err != nil {
		t.Fatal(err)
	}
	worktree, _ = filepath.EvalSymlinks(worktree)
	project := ComposeProjectName(name)

	sh := func(name string, args ...string) string {
		t.Helper()
		out, err := exec.Command(name, args...).CombinedOutput()
		if err != nil {
			t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
		}
		return strings.TrimSpace(string(out))
	}
	labels := []string{"--label", composeProjectLabel + "=" + project, "--label", composeDirLabel + "=" + worktree}
	sh("docker", append(append([]string{"network", "create", "--subnet", "10.254.254.0/29"}, labels...), project+"_default")...)
	sh("docker", append(append([]string{"volume", "create"}, labels...), project+"_data")...)
	sh("docker", append(append([]string{"run", "-d", "--name", project + "-app-1", "--network", project + "_default", "-v", project + "_data:/data"}, labels...), "redis:alpine")...)
	foreign := name + "-foreign"
	sh("docker", "run", "-d", "--name", foreign, "--label", composeProjectLabel+"=other-"+project, "--label", composeDirLabel+"=/elsewhere", "redis:alpine")
	t.Cleanup(func() { _ = exec.Command("docker", "rm", "-f", foreign).Run() })

	udid := sh("xcrun", "simctl", "create", "Probe-"+name, "com.apple.CoreSimulator.SimDeviceType.iPhone-17-Pro")
	t.Cleanup(func() { _ = exec.Command("xcrun", "simctl", "delete", udid).Run() })

	derived := t.TempDir()
	dd := filepath.Join(derived, "Probe-abcdef")
	if err := os.MkdirAll(dd, 0o755); err != nil {
		t.Fatal(err)
	}
	plist := `<?xml version="1.0" encoding="UTF-8"?><plist version="1.0"><dict><key>WorkspacePath</key><string>` + worktree + `/Probe.xcodeproj</string></dict></plist>`
	if err := os.WriteFile(filepath.Join(dd, "info.plist"), []byte(plist), 0o644); err != nil {
		t.Fatal(err)
	}

	sleeper := exec.Command("sleep", "600")
	sleeper.Dir = worktree
	if err := sleeper.Start(); err != nil {
		t.Fatal(err)
	}
	exited := make(chan struct{})
	go func() { _ = sleeper.Wait(); close(exited) }()

	c := Live()
	for i, s := range c.Sweeps {
		if d, ok := s.(DerivedData); ok {
			d.Root = derived
			c.Sweeps[i] = d
		}
	}
	r := c.Clean(ctx, Target{Worktree: worktree})
	t.Logf("cleanup: %s", r.Meta())
	if len(r.Pending) > 0 {
		t.Fatalf("pending: %q", r.Pending)
	}

	if out := sh("docker", "ps", "-aq", "--filter", "label="+composeProjectLabel+"="+project); out != "" {
		t.Errorf("containers left: %q", out)
	}
	if out := sh("docker", "volume", "ls", "-q", "--filter", "label="+composeProjectLabel+"="+project); out != "" {
		t.Errorf("volumes left: %q", out)
	}
	if out := sh("docker", "network", "ls", "-q", "--filter", "label="+composeProjectLabel+"="+project); out != "" {
		t.Errorf("networks left: %q", out)
	}
	if out := sh("docker", "ps", "-aq", "--filter", "name="+foreign); out == "" {
		t.Error("the foreign container was removed")
	}
	if strings.Contains(sh("xcrun", "simctl", "list", "devices"), udid) {
		t.Error("the simulator is still there")
	}
	if _, err := os.Stat(dd); !os.IsNotExist(err) {
		t.Errorf("DerivedData folder left: %v", err)
	}
	select {
	case <-exited:
	case <-time.After(10 * time.Second):
		_ = sleeper.Process.Kill()
		t.Error("the process in the worktree is still running")
	}
}
