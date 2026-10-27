package claude

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// A topic's session (CON-104) is a fork of the Mate's session: `--resume
// <mate> --fork-session` on its first run, `--resume <fork>` after. It
// must neither act nor pass for the Mate. Measured 2026-10-20 on Claude
// Code 2.1.292 (CON-104 comment):
//   - the Mate's Stop hook writes the run's session id into mate.meta, so a
//     fork that ran it would take the Mate's session over; the fork's
//     settings are the Mate's without hooks and with disableAllHooks;
//   - with --permission-mode dontAsk and no --dangerously-skip-permissions
//     every Bash command outside --allowedTools is denied, redirections
//     and `a && b` included, and --tools without Write or Edit leaves no
//     way to change a file;
//   - --strict-mcp-config without --mcp-config loads no MCP server, so the
//     plugins' write tools are gone too;
//   - dontAsk also denies Read, Glob and Grep outside the Mate's dir, so a
//     fork could read neither the project's repo nor a screenshot the
//     captain pasted (2026-10-27); ForkReadTools are allowed everywhere.

// ForkTools are the only tools a fork has.
const ForkTools = "Read,Glob,Grep,Bash"

// ForkReadTools are the fork's tools allowed on any path: they read only.
var ForkReadTools = []string{"Read", "Glob", "Grep"}

// ForkSettingsFile is the fork's settings file, in the dir ForkTurnArgs is
// given.
const ForkSettingsFile = "fork-settings.json"

// forkDropped are launch flags a fork run does not keep, with whether each
// takes a value.
var forkDropped = map[string]bool{
	"--dangerously-skip-permissions": false,
	"--resume":                       true,
	"--session-id":                   true,
	"--settings":                     true,
	"--tools":                        true,
	"--mcp-config":                   true,
	"--permission-mode":              true,
}

// ForkTurnArgs implements harness.ForkTurns.
func (Claude) ForkTurnArgs(launch []string, fork string, allow []string, dir string) ([]string, error) {
	return ForkTurnArgs(launch, fork, allow, dir)
}

// ForkTurnArgs is a fork run's argv: the Mate's launch args without what
// lets it act or names its session, the fork's settings, tools and
// permission rules, the session to fork or resume, then StreamTurnFlags.
func ForkTurnArgs(launch []string, fork string, allow []string, dir string) ([]string, error) {
	from := LaunchSessionID(launch)
	if from == "" && fork == "" {
		return nil, errors.New("claude: the launch names no session to fork")
	}
	settings := ""
	out := make([]string, 0, len(launch)+len(allow)+len(StreamTurnFlags)+10)
	for i := 0; i < len(launch); i++ {
		takesValue, drop := forkDropped[launch[i]]
		if !drop {
			out = append(out, launch[i])
			continue
		}
		if takesValue && i+1 < len(launch) {
			if launch[i] == "--settings" {
				settings = launch[i+1]
			}
			i++
		}
	}
	path, err := writeForkSettings(settings, dir)
	if err != nil {
		return nil, err
	}
	out = append(out, "--settings", path, "--tools", ForkTools, "--permission-mode", "dontAsk", "--allowedTools")
	out = append(out, ForkReadTools...)
	out = append(out, allow...)
	if fork == "" {
		out = append(out, "--resume", from, "--fork-session")
	} else {
		out = append(out, "--resume", fork)
	}
	return append(out, StreamTurnFlags...), nil
}

// writeForkSettings writes the Mate's settings (src, none when empty)
// without its hooks and with every hook off into dir, and returns the
// file's path.
func writeForkSettings(src, dir string) (string, error) {
	settings := map[string]any{}
	if src != "" {
		data, err := os.ReadFile(src)
		if err != nil {
			return "", fmt.Errorf("claude: fork settings: %w", err)
		}
		if err := json.Unmarshal(data, &settings); err != nil {
			return "", fmt.Errorf("claude: fork settings: %s: %w", src, err)
		}
	}
	delete(settings, "hooks")
	settings["disableAllHooks"] = true
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	path := filepath.Join(dir, ForkSettingsFile)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o600); err != nil {
		return "", err
	}
	return path, os.Rename(tmp, path)
}
