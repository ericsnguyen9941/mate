package pi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	if err := checkGlobalMCP(servers); err != nil {
		return nil, err
	}
	return json.MarshalIndent(map[string]any{"mcpServers": entries}, "", "  ")
}

// Global entries take precedence over registered servers in pi. Refuse a
// collision instead of launching a workspace instance with a different URL.
func checkGlobalMCP(servers []harness.MCPServer) error {
	if len(servers) == 0 {
		return nil
	}
	root, err := sessionsRoot("")
	if err != nil {
		return err
	}
	path := filepath.Join(filepath.Dir(root), "mcp.json")
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var doc struct {
		Servers map[string]json.RawMessage `json:"mcpServers"`
	}
	if err = json.Unmarshal(b, &doc); err != nil {
		return fmt.Errorf("pi global MCP config %s: %w", path, err)
	}
	for _, s := range servers {
		for n := range doc.Servers {
			if strings.ReplaceAll(n, "-", "_") == strings.ReplaceAll(s.Name, "-", "_") {
				return observability.NewError(observability.CodeUsage, fmt.Sprintf("pi MCP server %q conflicts with global entry %q in %s; rename the workspace instance or global entry", s.Name, n, path))
			}
		}
	}
	return nil
}
