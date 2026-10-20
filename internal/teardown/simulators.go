package teardown

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/process"
)

// Simulators deletes the iOS simulators named after the worktree. Nothing
// links a simulator to a directory, so the name is the convention: a repo
// that clones a simulator per worktree puts the worktree's directory name
// in it (plantrip's `PlanTrip-<worktree>`). A base device and a hand-made
// one carry no worktree name and are never matched.
type Simulators struct {
	Runner process.Runner
}

type simDevice struct {
	Name  string `json:"name"`
	UDID  string `json:"udid"`
	State string `json:"state"`
	Size  int64  `json:"dataPathSize"`
}

// Sweep implements Sweep.
func (s Simulators) Sweep(ctx context.Context, targets []Target, dry bool) ([]string, error) {
	res, err := s.run(ctx, "list", "devices", "-j")
	if err != nil {
		if missing(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("simulators: %v", err)
	}
	if res.ExitCode != 0 {
		return nil, fmt.Errorf("simulators: simctl list exited %d%s", res.ExitCode, lastLine(res.Stderr))
	}
	var list struct {
		Devices map[string][]simDevice `json:"devices"`
	}
	if err := json.Unmarshal(res.Stdout, &list); err != nil {
		return nil, fmt.Errorf("simulators: simctl list: %v", err)
	}
	var owned []simDevice
	for _, devices := range list.Devices {
		for _, d := range devices {
			if d.UDID != "" && namedAfter(targets, d.Name) {
				owned = append(owned, d)
			}
		}
	}
	sort.Slice(owned, func(i, j int) bool { return owned[i].Name < owned[j].Name })

	var removed, failed []string
	for _, d := range owned {
		if dry {
			removed = append(removed, "simulator "+d.Name+sizeNote(d.Size))
			continue
		}
		if d.State != "Shutdown" {
			// A device that is already shut down answers with an error;
			// only the delete below decides.
			_, _ = s.run(ctx, "shutdown", d.UDID)
		}
		res, err := s.run(ctx, "delete", d.UDID)
		switch {
		case err != nil:
			failed = append(failed, fmt.Sprintf("simulator %s: %v", d.Name, err))
		case res.ExitCode != 0:
			failed = append(failed, fmt.Sprintf("simulator %s: delete exited %d%s", d.Name, res.ExitCode, lastLine(res.Stderr)))
		default:
			removed = append(removed, "simulator "+d.Name+sizeNote(d.Size))
		}
	}
	if len(failed) > 0 {
		return removed, fmt.Errorf("%s", strings.Join(failed, "; "))
	}
	return removed, nil
}

func (s Simulators) run(ctx context.Context, args ...string) (process.Result, error) {
	run := s.Runner
	if run == nil {
		run = process.ExecRunner{}
	}
	return run.Run(ctx, process.Spec{Name: "xcrun", Args: append([]string{"simctl"}, args...)})
}
