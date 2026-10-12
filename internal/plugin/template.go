package plugin

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

var (
	configVarRe = regexp.MustCompile(`^config\.[A-Za-z0-9_-]+$`)
	// skillVarRe matches what a skill file treats as a variable reference.
	// Other braces, such as JSON in code samples, are left as written.
	skillVarRe = regexp.MustCompile(`\{[A-Za-z_][A-Za-z0-9_.]*\}`)
)

// expandURL substitutes {config.<key>} literally from values and checks the
// result is an http(s) URL.
func expandURL(tmpl string, values map[string]string) (string, error) {
	u, _, err := substitute(tmpl, func(key string) (string, bool) {
		v, ok := values[key]
		return v, ok
	})
	if err != nil {
		return "", err
	}
	if err := checkURL(u); err != nil {
		return "", err
	}
	return u, nil
}

// substitute replaces each {config.<key>} in tmpl with lookup(key) and
// reports how many it replaced. Every brace must open such a variable: there
// is no expression, escaping, or recursive expansion.
func substitute(tmpl string, lookup func(string) (string, bool)) (string, int, error) {
	var out strings.Builder
	s, n := tmpl, 0
	for {
		i := strings.IndexAny(s, "{}")
		if i < 0 {
			out.WriteString(s)
			return out.String(), n, nil
		}
		if s[i] == '}' {
			return "", 0, fmt.Errorf("unmatched } in %q", tmpl)
		}
		out.WriteString(s[:i])
		s = s[i+1:]
		end := strings.IndexByte(s, '}')
		if end < 0 {
			return "", 0, fmt.Errorf("unclosed variable in %q", tmpl)
		}
		name := s[:end]
		if !configVarRe.MatchString(name) {
			return "", 0, fmt.Errorf("unknown variable {%s} in %q; only {config.<key>} is supported", name, tmpl)
		}
		key := strings.TrimPrefix(name, "config.")
		value, ok := lookup(key)
		if !ok {
			return "", 0, fmt.Errorf("%q needs {config.%s}, but the binding config has no %s", tmpl, key, key)
		}
		out.WriteString(value)
		n++
		s = s[end+1:]
	}
}

// checkURL accepts only absolute http and https URLs with a host and no
// userinfo, so no credentials live in plugin config or resolved URLs.
func checkURL(s string) error {
	if strings.IndexFunc(s, func(r rune) bool { return unicode.IsSpace(r) || unicode.IsControl(r) }) >= 0 {
		return fmt.Errorf("URL %q contains whitespace or control characters", s)
	}
	u, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid URL %q: %v", s, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL %q must use http or https", s)
	}
	if u.Host == "" {
		return fmt.Errorf("URL %q has no host", s)
	}
	if u.User != nil {
		return fmt.Errorf("URL %q must not contain credentials", s)
	}
	return nil
}

// expandSkill substitutes {instance_id}; any other variable reference is an
// error so a skill never ships with a placeholder Mate does not fill. Shell
// expansions such as ${HOME} are not variable references.
func expandSkill(content []byte, instance string) ([]byte, error) {
	var out []byte
	last := 0
	for _, loc := range skillVarRe.FindAllIndex(content, -1) {
		if loc[0] > 0 && content[loc[0]-1] == '$' {
			continue
		}
		name := content[loc[0]:loc[1]]
		if string(name) != "{instance_id}" {
			return nil, fmt.Errorf("unknown variable %s; only {instance_id} is supported", name)
		}
		out = append(out, content[last:loc[0]]...)
		out = append(out, instance...)
		last = loc[1]
	}
	return append(out, content[last:]...), nil
}
