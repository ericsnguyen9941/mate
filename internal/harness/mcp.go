package harness

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/observability"
)

// ValidateMCPServers is the common transport contract, used before launch.
func ValidateMCPServers(servers []MCPServer) error {
	seen := map[string]bool{}
	for _, s := range servers {
		fail := func(msg string) error {
			return observability.NewError(observability.CodeUsage, fmt.Sprintf("MCP server %q: %s", s.Name, msg))
		}
		if strings.TrimSpace(s.Name) == "" || strings.ContainsAny(s.Name, "\x00\r\n") {
			return fail("invalid name")
		}
		if seen[s.Name] {
			return fail("duplicate name")
		}
		seen[s.Name] = true
		switch s.Type {
		case "http":
			u, err := url.Parse(s.URL)
			if err != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") || u.User != nil || len(s.Command) > 0 || s.Cwd != "" {
				return fail("http needs an http(s) URL without userinfo and takes no command or cwd")
			}
		case "stdio":
			if len(s.Command) == 0 || s.Command[0] == "" || s.URL != "" {
				return fail("stdio needs a command and takes no URL")
			}
			for _, a := range s.Command {
				if strings.ContainsRune(a, 0) {
					return fail("command contains NUL")
				}
			}
			if s.Cwd != "" && (!filepath.IsAbs(s.Cwd) || strings.ContainsRune(s.Cwd, 0)) {
				return fail("cwd must be absolute")
			}
		default:
			return fail("transport must be http or stdio")
		}
	}
	return nil
}
