package claude

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/observability"
)

// MCPConfigFile is the --mcp-config file of a Claude Mate whose project has
// plugin MCP servers, in the Mate's cwd. Mate owns it: every Prepare
// rewrites it from PrepareRequest.MCPServers.
const MCPConfigFile = "mcp-plugins.json"

// MCPConfigPath is the --mcp-config file of a Claude Mate in mateDir.
func MCPConfigPath(mateDir string) string { return filepath.Join(mateDir, MCPConfigFile) }

// mcpCwdShell starts a stdio server in its own directory. Claude Code
// 2.1.287 has no cwd for a stdio server: its config schema takes type,
// command, args, env, timeout and alwaysLoad, an unknown "cwd" key is
// dropped, and the server process runs in the session's cwd. Measured
// 2026-10-12: a server given `"cwd": "<pkg>"` with `"command":
// "./mcpfixture"` failed to start, and `pwd -P` from a server given that
// cwd printed the session's cwd (docs/evidence/claude-mcp-probe-2026-10-12.md,
// "PR 4 verification"). So the server's own directory is reached the way a
// shell would: `cd` into it, then exec the command. The directory and the
// command are positional parameters, never part of the script.
var mcpCwdShell = []string{"/bin/sh", "-c", `cd "$1" && shift && exec "$@"`, "mate-mcp-cwd"}

// mcpConfig is Claude Code's --mcp-config document: `{"mcpServers":
// {"<name>": <server>}}`, the shape `claude mcp add-json` takes and the
// 2.1.287 schema reads (type "http" with url; stdio with command and args,
// type optional). Nothing else goes in it: no env, no header, no token.
// Auth is Claude's own (its OAuth for an http server), never Mate's.
type mcpConfig struct {
	MCPServers map[string]mcpServerEntry `json:"mcpServers"`
}

// mcpServerEntry is one server: Type and URL for http, Command and Args for
// stdio. Fields of the other transport are left out.
type mcpServerEntry struct {
	Type    string    `json:"type,omitempty"`
	URL     string    `json:"url,omitempty"`
	Command string    `json:"command,omitempty"`
	Args    *[]string `json:"args,omitempty"`
}

// MCPConfig is the --mcp-config document for servers. A server Claude
// cannot be given as it is described is refused, not dropped.
func MCPConfig(servers []harness.MCPServer) ([]byte, error) {
	doc := mcpConfig{MCPServers: map[string]mcpServerEntry{}}
	for _, s := range servers {
		if s.Name == "" {
			return nil, mcpRefusal("an MCP server has no name")
		}
		if _, dup := doc.MCPServers[s.Name]; dup {
			return nil, mcpRefusal(fmt.Sprintf("two MCP servers are named %q", s.Name))
		}
		entry, err := mcpEntry(s)
		if err != nil {
			return nil, err
		}
		doc.MCPServers[s.Name] = entry
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doc); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func mcpEntry(s harness.MCPServer) (mcpServerEntry, error) {
	switch s.Type {
	case "http":
		if s.URL == "" || len(s.Command) > 0 || s.Cwd != "" {
			return mcpServerEntry{}, mcpRefusal(fmt.Sprintf("MCP server %q is http: it needs a url and takes no command or cwd", s.Name))
		}
		return mcpServerEntry{Type: "http", URL: s.URL}, nil
	case "stdio":
		if len(s.Command) == 0 || s.Command[0] == "" || s.URL != "" {
			return mcpServerEntry{}, mcpRefusal(fmt.Sprintf("MCP server %q is stdio: it needs a command and takes no url", s.Name))
		}
		argv := s.Command
		if s.Cwd != "" {
			if !filepath.IsAbs(s.Cwd) {
				return mcpServerEntry{}, mcpRefusal(fmt.Sprintf("MCP server %q has a relative cwd %q", s.Name, s.Cwd))
			}
			argv = append(append(append([]string(nil), mcpCwdShell...), s.Cwd), s.Command...)
		}
		args := append([]string{}, argv[1:]...)
		return mcpServerEntry{Command: argv[0], Args: &args}, nil
	default:
		return mcpServerEntry{}, mcpRefusal(fmt.Sprintf("MCP server %q has type %q; Claude is given http or stdio", s.Name, s.Type))
	}
}

func mcpRefusal(msg string) error {
	return observability.NewError(observability.CodeUsage, "claude launch: "+msg)
}

// mateMCPFile is the --mcp-config file of a Mate launch and its path, which
// Build passes on; both are empty when the launch has no servers. A launch
// with none leaves no earlier list behind: a file still there from an
// earlier start is rewritten with no servers. It is not removed because a
// Prepare only names files for spawn to write (harness.Prepared), and with
// no --mcp-config on the argv Claude never reads it. Only when the file
// exists is it named, so a Mate without plugins has the files it always had.
func mateMCPFile(req harness.PrepareRequest) (*harness.LaunchFile, string, error) {
	path := MCPConfigPath(req.Cwd)
	if len(req.MCPServers) == 0 {
		if _, err := os.Lstat(path); err != nil {
			if os.IsNotExist(err) {
				return nil, "", nil
			}
			return nil, "", err
		}
		data, err := MCPConfig(nil)
		if err != nil {
			return nil, "", err
		}
		return &harness.LaunchFile{Path: path, Data: data}, "", nil
	}
	data, err := MCPConfig(req.MCPServers)
	if err != nil {
		return nil, "", err
	}
	return &harness.LaunchFile{Path: path, Data: data}, path, nil
}
