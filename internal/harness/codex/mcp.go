package codex

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

type codexMCP struct{}

func (codexMCP) MCPConfig(servers []harness.MCPServer) ([]byte, error) { return MCPConfig(servers) }

// MCPConfig renders a TOML inline table, passed as one root-table override.
// Replacing the whole table prevents global entries or old transport fields
// from merging into a workspace instance. It never writes config.toml.
func MCPConfig(servers []harness.MCPServer) ([]byte, error) {
	if err := harness.ValidateMCPServers(servers); err != nil {
		return nil, err
	}
	sorted := append([]harness.MCPServer(nil), servers...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })
	quote := func(v any) string { b, _ := json.Marshal(v); return string(b) }
	entries := []string{}
	for _, s := range sorted {
		fields := []string{}
		if s.Type == "http" {
			fields = append(fields, "url="+quote(s.URL))
		} else {
			fields = append(fields, "command="+quote(s.Command[0]), "args="+quote(append([]string{}, s.Command[1:]...)))
			if s.Cwd != "" {
				fields = append(fields, "cwd="+quote(s.Cwd))
			}
		}
		fields = append(fields, "enabled=true")
		entries = append(entries, quote(s.Name)+"={"+strings.Join(fields, ",")+"}")
	}
	return []byte("{" + strings.Join(entries, ",") + "}"), nil
}

// MCPArgs can also be used with `codex mcp login <name>` in the same home.
func MCPArgs(servers []harness.MCPServer) ([]string, error) {
	doc, err := MCPConfig(servers)
	if err != nil {
		return nil, err
	}
	return []string{"-c", "mcp_servers=" + string(doc)}, nil
}
