package spawn

import (
	"errors"
	"fmt"

	"github.com/nguyenngocanh94/mate/internal/harness"
	"github.com/nguyenngocanh94/mate/internal/mateassets"
	"github.com/nguyenngocanh94/mate/internal/observability"
	"github.com/nguyenngocanh94/mate/internal/plugin"
	"github.com/nguyenngocanh94/mate/internal/store"
)

// matePlugins is what the project's Mate receives from the workspace's
// plugins.yaml (internal/plugin): the MCP servers its launch exposes and the
// skills written beside its own.
type matePlugins struct {
	servers []harness.MCPServer
	skills  []mateassets.PluginSkill
}

// resolveMatePlugins reads the project's plugins for a Mate start. Only a
// Mate start asks; a Crew gets no plugin. A configuration problem refuses
// the start with the file, instance and project it is in, rather than
// starting a Mate without its plugins. The Mate picks up a change on its
// next start; nothing restarts it.
func resolveMatePlugins(w *store.Workspace, project string) (matePlugins, error) {
	r, err := plugin.Resolve(w.Root(), project)
	if err != nil {
		return matePlugins{}, observability.WrapError(observability.CodeUsage,
			fmt.Sprintf("the Mate of %s cannot start: its plugin configuration is invalid", project), err)
	}
	var p matePlugins
	for _, s := range r.Servers {
		p.servers = append(p.servers, harness.MCPServer{
			Name: s.Name, Type: s.Type, URL: s.URL,
			Command: append([]string(nil), s.Command...), Cwd: s.Cwd,
		})
	}
	for _, s := range r.Skills {
		p.skills = append(p.skills, mateassets.PluginSkill{Instance: s.Instance, Name: s.Name, Content: s.Content})
	}
	return p, nil
}

// PreflightMatePlugins says whether a Mate start of project on kind (the
// workspace default when empty) would get past its plugins: their
// configuration resolves, and a harness whose MCPConfig capability is not
// verified has no MCP server to be given. It is the refusal StartMate or
// the harness's Prepare would return, asked before anything is stopped:
// context refresh asks it so that a Mate is not stopped for a fresh start
// that cannot come up. It reads files only.
func PreflightMatePlugins(w *store.Workspace, deps Deps, project string, kind harness.Kind) error {
	if kind == "" {
		var err error
		if kind, err = deps.defaultHarness(w.Defaults().MateHarness, harness.RoleMate); err != nil {
			return err
		}
	}
	profile, err := deps.Harnesses.Lookup(kind)
	if err != nil {
		return err
	}
	p, err := resolveMatePlugins(w, project)
	if err != nil {
		return err
	}
	if c := profile.Capabilities().MCPConfig; !c.Verified() {
		return harness.RefuseMCPServers(kind, c.Reason, p.servers)
	}
	return nil
}

// writePluginSkills installs the plugin skills where the harness reads the
// Mate's own, and removes those of instances the Mate no longer has. A
// destination the user owns refuses the start.
func writePluginSkills(mateDir string, profile harness.Profile, skills []mateassets.PluginSkill) error {
	err := mateassets.WritePluginSkills(mateDir, profile.Info().SkillsDir, skills)
	if errors.Is(err, mateassets.ErrPluginSkillConflict) {
		return observability.WrapError(observability.CodeStateConflict, "a plugin skill would overwrite a skill mate did not write", err)
	}
	return err
}
