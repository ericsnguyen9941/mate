package plugin

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestWorkspaceInstallationReachesCurrentAndFutureProjects(t *testing.T) {
	root := emptyWorkspace(t)
	_, err := AddBinding(root, AddRequest{Instance: "tracker-work", Path: filepath.Join(packagesDir(t), "full"), Config: map[string]string{"project_url": "https://tracker.example/workspace"}, Key: "t"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(readPlugins(t, root), "projects:") {
		t.Fatal("installation wrote project bindings")
	}
	var first Resolved
	for _, project := range []string{"shop", "blog", "created-later"} {
		got, err := Resolve(root, project)
		if err != nil {
			t.Fatal(err)
		}
		if len(got.Servers) != 1 || len(got.Skills) != 1 || len(got.Links) != 1 || got.Links[0].Key != "t" {
			t.Fatalf("%s: %+v", project, got)
		}
		if project == "shop" {
			first = got
		} else if !reflect.DeepEqual(first, got) {
			t.Fatalf("%s does not inherit workspace plugin", project)
		}
	}
	checks := Doctor(root, DoctorOptions{Projects: []string{"shop", "blog"}, Harness: func(string) (string, string) { return "claude", "workspace default" }})
	for _, check := range checks {
		if !check.OK() {
			t.Fatal(check.String())
		}
	}
}

func TestWorkspacePromotionPreservesOtherInstances(t *testing.T) {
	root := workspace(t, `apiVersion: mate.plugins/v1
instances:
  notes:
    path: $PKG/skillsonly
    projects:
      shop: {}
      blog: {}
  other:
    path: $PKG/skillsonly
    projects:
      blog: {}
`)
	_, err := AddBinding(root, AddRequest{Instance: "notes", Path: filepath.Join(packagesDir(t), "skillsonly")})
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Instances[0].Workspace == nil || len(cfg.Instances[0].Projects) != 0 || len(cfg.Instances[1].Projects) != 1 {
		t.Fatalf("promotion: %+v", cfg.Instances)
	}
	got, err := Resolve(root, "future")
	if err != nil || len(got.Skills) != 2 {
		t.Fatalf("skills-only workspace plugin: %+v, %v", got, err)
	}
}

func TestWorkspaceConfigurationValidation(t *testing.T) {
	base := `apiVersion: mate.plugins/v1
instances:
  shared:
    path: $PKG/full
    config: {project_url: "https://tracker.example/workspace"}
    key: t
`
	for _, tc := range []struct{ name, extra, project, want string }{
		{"workspace collision", "  another:\n    path: $PKG/full\n    config: {project_url: 'https://tracker.example/other'}\n    key: t\n", "future", "duplicate key"},
		{"legacy collision", "  legacy:\n    path: $PKG/full\n    projects:\n      shop:\n        config: {project_url: 'https://tracker.example/shop'}\n        key: t\n", "shop", "duplicate key"},
		{"malformed workspace entry", "  broken: []\n", "future", "cannot unmarshal"},
		{"malformed workspace", "  broken:\n    path: $PKG/skillsonly\n    typo: true\n", "future", "field typo"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Resolve(workspace(t, base+tc.extra), tc.project)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want %q, got %v", tc.want, err)
			}
		})
	}
}
