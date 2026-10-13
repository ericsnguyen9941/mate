package pi

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/observability"
)

type piMCP struct{}

func (piMCP) MCPConfig(servers []harness.MCPServer) ([]byte, error) { return MCPConfig(servers) }

var serverName = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// MCPConfig uses pi's native registered-server schema. Direct exposure makes
// the first prompt wait for discovery and supplies the normal MCP tool names.
func MCPConfig(servers []harness.MCPServer) ([]byte, error) {
	if err := harness.ValidateMCPServers(servers); err != nil {
		return nil, err
	}
	entries := map[string]map[string]any{}
	names := map[string]string{}
	for _, s := range servers {
		if !serverName.MatchString(s.Name) {
			return nil, observability.NewError(observability.CodeUsage, fmt.Sprintf("pi MCP server %q: names require letters, digits, _ or -", s.Name))
		}
		ns := strings.ReplaceAll(s.Name, "-", "_")
		if other, ok := names[ns]; ok {
			return nil, observability.NewError(observability.CodeUsage, fmt.Sprintf("pi MCP names %q and %q share a tool namespace", other, s.Name))
		}
		names[ns] = s.Name
		e := map[string]any{"exposure": "direct"}
		if s.Type == "http" {
			e["url"] = s.URL
		} else {
			e["command"] = s.Command[0]
			e["args"] = append([]string{}, s.Command[1:]...)
			if s.Cwd != "" {
				e["cwd"] = s.Cwd
			}
		}
		entries[s.Name] = e
	}
	return json.MarshalIndent(map[string]any{"mcpServers": entries}, "", "  ")
}
