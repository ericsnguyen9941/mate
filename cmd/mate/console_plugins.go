package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"time"

	"github.com/nguyenngocanh94/mate/internal/consolecmd"
	"github.com/nguyenngocanh94/mate/internal/plugin"
	"github.com/nguyenngocanh94/mate/internal/query"
	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/ui/console"
)

// pluginLinks are the plugin links that carry a console key, by project.
type pluginLinks map[string][]plugin.Link

// loadPluginLinks resolves every project's plugin links. A plugin
// configuration error leaves the console without any plugin link and is
// returned as notice, one line naming the file and the error, for the
// status line: plugins are a convenience the console runs without. A key
// that console.yaml or the console already owns is console configuration,
// and stops the console before it starts; nothing is overridden silently.
func loadPluginLinks(ws *store.Workspace, cfg consolecmd.Config) (links pluginLinks, notice string, err error) {
	links = pluginLinks{}
	for _, p := range ws.Projects() {
		res, err := plugin.Resolve(ws.Root(), p.Name)
		if err != nil {
			return pluginLinks{}, "plugin links off: " + strings.Join(strings.Fields(err.Error()), " "), nil
		}
		for _, l := range res.Links {
			if l.Key == "" {
				continue
			}
			if owner := consoleKeyOwner(cfg, l.Key); owner != "" {
				return nil, "", fmt.Errorf("%s: instance %q key %q is %s", plugin.ConfigFile, l.Instance, l.Key, owner)
			}
			links[p.Name] = append(links[p.Name], l)
		}
	}
	return links, "", nil
}

// consoleKeyOwner says what already owns key on project rows - the console
// itself or a project binding in console.yaml - or "" when a plugin link
// may take it. The console start and `mate plugin` ask the same question.
func consoleKeyOwner(cfg consolecmd.Config, key string) string {
	if slices.Contains(consolecmd.ReservedKeys(), key) {
		return "reserved by the console"
	}
	for _, b := range cfg.Bindings {
		if b.Scope == "project" && b.Key == key {
			return "already bound on project rows in .mate/console.yaml"
		}
	}
	return ""
}

// bindings adds one key-line entry per key and instance, however many
// projects bind it.
func (l pluginLinks) bindings() []query.CommandBinding {
	var out []query.CommandBinding
	seen := map[string]bool{}
	for _, p := range sortedProjects(l) {
		for _, link := range l[p] {
			id := link.Key + "\x00" + link.Instance
			if seen[id] {
				continue
			}
			seen[id] = true
			out = append(out, query.CommandBinding{Key: link.Key, Label: link.Title, Scope: "project", Role: "plugin-" + link.Instance})
		}
	}
	return out
}

func sortedProjects(l pluginLinks) []string {
	names := make([]string, 0, len(l))
	for p := range l {
		names = append(names, p)
	}
	slices.Sort(names)
	return names
}

func (l pluginLinks) hasKey(key string) bool {
	for _, links := range l {
		for _, link := range links {
			if link.Key == key {
				return true
			}
		}
	}
	return false
}

// urlOpener names the machine's URL opener.
func urlOpener(goos string) string {
	switch goos {
	case "darwin":
		return "open"
	case "linux":
		return "xdg-open"
	}
	return ""
}

// consolePluginView opens plugin link URLs with the machine's opener and
// hands every other key to inner. No terminal tab is involved.
func consolePluginView(links pluginLinks, opener string, runner consolecmd.Runner, inner console.CommandViewFunc) console.CommandViewFunc {
	if len(links) == 0 {
		return inner
	}
	return func(ctx context.Context, key string, target console.StageTarget) error {
		if target.Kind == "" && target.ID == "" && links.hasKey(key) {
			for _, l := range links[target.ProjectID] {
				if l.Key == key {
					return openURL(ctx, opener, runner, l)
				}
			}
			return fmt.Errorf("no app bound to %s for project %s", key, target.ProjectID)
		}
		if inner == nil {
			return fmt.Errorf("no command bound to %s on this row", key)
		}
		return inner(ctx, key, target)
	}
}

func openURL(ctx context.Context, opener string, runner consolecmd.Runner, l plugin.Link) error {
	name := opener
	if name == "" {
		return fmt.Errorf("%s: no URL opener for %s", l.Title, runtime.GOOS)
	}
	if !filepath.IsAbs(name) {
		name = findTool(os.Getenv, name)
	}
	if !filepath.IsAbs(name) {
		return fmt.Errorf("%s: URL opener %q not found", l.Title, opener)
	}
	launchCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var diagnostic bytes.Buffer
	if err := runner(launchCtx, consolecmd.Invocation{Name: name, Args: []string{l.URL}}, nil, &diagnostic, &diagnostic); err != nil {
		return fmt.Errorf("%s: %w: %s", l.Title, err, strings.TrimSpace(diagnostic.String()))
	}
	return nil
}

func hostOpener() string { return urlOpener(runtime.GOOS) }
