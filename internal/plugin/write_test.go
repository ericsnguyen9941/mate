package plugin

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func emptyWorkspace(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".mate"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

func readPlugins(t *testing.T, root string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, ".mate", "plugins.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestAddBindingCreatesTheFile(t *testing.T) {
	pkgs := packagesDir(t)
	root := emptyWorkspace(t)
	res, err := AddBinding(root, AddRequest{
		Instance: "tracker-work", Path: filepath.Join(pkgs, "full"), Project: "shop",
		Config: map[string]string{"project_url": "https://tracker.example/p/shop"}, Key: "t",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Created || !res.NewInstance || res.Replaced || res.Disabled || res.Manifest == nil || res.Manifest.Name != "Tracker" {
		t.Fatalf("result = %+v", res)
	}
	want := "apiVersion: mate.plugins/v1\ninstances:\n  tracker-work:\n    path: " + filepath.Join(pkgs, "full") +
		"\n    projects:\n      shop:\n        config:\n          project_url: https://tracker.example/p/shop\n        key: t\n"
	if got := readPlugins(t, root); got != want {
		t.Fatalf("plugins.yaml =\n%s\nwant\n%s", got, want)
	}
	r, err := Resolve(root, "shop")
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Links) != 1 || r.Links[0].URL != "https://tracker.example/p/shop" || r.Links[0].Key != "t" {
		t.Fatalf("links = %+v", r.Links)
	}
}

func TestAddBindingKeepsOtherInstances(t *testing.T) {
	pkgs := packagesDir(t)
	other := `# workspace plugins
apiVersion: mate.plugins/v1
instances:
  # the notes app
  notes:
    path: $PKG/skillsonly # kept
    projects:
      blog: {}
      shop:
        config:
          id: "42"
  board:
    path: $PKG/openonly
    enabled: false
    projects:
      blog:
        config:
          team: core
          project_id: "7"
`
	root := workspace(t, other)
	before := readPlugins(t, root)
	if _, err := AddBinding(root, AddRequest{Instance: "local", Path: filepath.Join(pkgs, "stdio"), Project: "shop"}); err != nil {
		t.Fatal(err)
	}
	got := readPlugins(t, root)
	want := before + "  local:\n    path: " + filepath.Join(pkgs, "stdio") + "\n    projects:\n      shop: {}\n"
	if got != want {
		t.Fatalf("plugins.yaml =\n%s\nwant\n%s", got, want)
	}
}

func TestAddBindingSameInstanceNewAndReplacedProject(t *testing.T) {
	pkgs := packagesDir(t)
	root := workspace(t, `apiVersion: mate.plugins/v1
instances:
  board:
    path: $PKG/openonly
    projects:
      blog:
        config:
          team: core
          project_id: "7"
        key: b
`)
	res, err := AddBinding(root, AddRequest{Instance: "board", Path: filepath.Join(pkgs, "openonly"), Project: "shop",
		Config: map[string]string{"team": "web", "project_id": "9"}})
	if err != nil {
		t.Fatal(err)
	}
	if res.Created || res.NewInstance || res.Replaced {
		t.Fatalf("result = %+v", res)
	}
	// Overwrite blog: the binding becomes exactly the request, key dropped.
	res, err = AddBinding(root, AddRequest{Instance: "board", Path: filepath.Join(pkgs, "openonly"), Project: "blog",
		Config: map[string]string{"team": "core", "project_id": "8"}})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Replaced {
		t.Fatalf("result = %+v, want Replaced", res)
	}
	want := `apiVersion: mate.plugins/v1
instances:
  board:
    path: $PKG/openonly
    projects:
      blog:
        config:
          project_id: "8"
          team: core
      shop:
        config:
          project_id: "9"
          team: web
`
	if got := readPlugins(t, root); got != strings.ReplaceAll(want, "$PKG", pkgs) {
		t.Fatalf("plugins.yaml =\n%s\nwant\n%s", got, strings.ReplaceAll(want, "$PKG", pkgs))
	}
}

func TestAddBindingRelativePathMatchesExistingInstance(t *testing.T) {
	pkgs := packagesDir(t)
	root := emptyWorkspace(t)
	if err := os.Symlink(pkgs, filepath.Join(root, "pkgs")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".mate", "plugins.yaml"),
		[]byte("apiVersion: mate.plugins/v1\ninstances:\n  notes:\n    path: pkgs/skillsonly\n    projects:\n      blog:\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := AddBinding(root, AddRequest{Instance: "notes", Path: filepath.Join(pkgs, "skillsonly"), Project: "shop"}); err != nil {
		t.Fatal(err)
	}
	got := readPlugins(t, root)
	if !strings.Contains(got, "    path: pkgs/skillsonly\n") || !strings.Contains(got, "      shop: {}\n") {
		t.Fatalf("plugins.yaml =\n%s", got)
	}
	if st, _ := os.Stat(filepath.Join(root, ".mate", "plugins.yaml")); st.Mode().Perm() != 0o600 {
		t.Fatalf("mode = %v, want the file's own 0600", st.Mode().Perm())
	}
}

func TestAddBindingRefusesWithoutWriting(t *testing.T) {
	pkgs := packagesDir(t)
	existing := `apiVersion: mate.plugins/v1
instances:
  board:
    path: $PKG/openonly
    projects:
      blog:
        config: {team: core, project_id: "7"}
        key: b
`
	for _, tc := range []struct {
		name string
		req  AddRequest
		want string
	}{{
		name: "different path",
		req:  AddRequest{Instance: "board", Path: filepath.Join(pkgs, "full"), Project: "shop"},
		want: "already uses package " + filepath.Join(pkgs, "openonly") + ", not " + filepath.Join(pkgs, "full") + "; to move it, change instances.board.path",
	}, {
		name: "missing config key",
		req:  AddRequest{Instance: "board", Path: filepath.Join(pkgs, "openonly"), Project: "shop", Config: map[string]string{"team": "web"}},
		want: "needs config project_id",
	}, {
		name: "missing config key on a new instance",
		req:  AddRequest{Instance: "tracker-work", Path: filepath.Join(pkgs, "full"), Project: "shop"},
		want: "needs config project_url",
	}, {
		name: "key without open link",
		req:  AddRequest{Instance: "notes", Path: filepath.Join(pkgs, "skillsonly"), Project: "shop", Key: "n"},
		want: `key "n" needs an open link`,
	}, {
		name: "duplicate key in project",
		req:  AddRequest{Instance: "tracker-work", Path: filepath.Join(pkgs, "full"), Project: "blog", Key: "b", Config: map[string]string{"project_url": "https://t.example/x"}},
		want: `duplicate key "b"`,
	}, {
		name: "package without manifest",
		req:  AddRequest{Instance: "nothing", Path: pkgs, Project: "shop"},
		want: "package has no mate-plugin.yaml",
	}, {
		name: "missing package",
		req:  AddRequest{Instance: "nothing", Path: filepath.Join(pkgs, "no-such"), Project: "shop"},
		want: "is not a directory",
	}, {
		name: "bad instance id",
		req:  AddRequest{Instance: "Board", Path: filepath.Join(pkgs, "openonly"), Project: "shop"},
		want: "instance id must match",
	}, {
		name: "invalid URL after substitution",
		req:  AddRequest{Instance: "tracker-work", Path: filepath.Join(pkgs, "full"), Project: "shop", Config: map[string]string{"project_url": "ftp://x"}},
		want: "must use http or https",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			root := workspace(t, existing)
			before := readPlugins(t, root)
			_, err := AddBinding(root, tc.req)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("err = %v, want %q", err, tc.want)
			}
			if got := readPlugins(t, root); got != before {
				t.Fatalf("plugins.yaml changed:\n%s", got)
			}
			entries, _ := os.ReadDir(filepath.Join(root, ".mate"))
			if len(entries) != 1 {
				t.Fatalf(".mate has %d entries, want plugins.yaml alone", len(entries))
			}
		})
	}
}

func TestAddBindingMissingConfigErrorNamesKeys(t *testing.T) {
	pkgs := packagesDir(t)
	root := emptyWorkspace(t)
	_, err := AddBinding(root, AddRequest{Instance: "board", Path: filepath.Join(pkgs, "openonly"), Project: "shop"})
	var me *MissingConfigError
	if !errors.As(err, &me) || !reflect.DeepEqual(me.Keys, []string{"project_id", "team"}) {
		t.Fatalf("err = %#v, want MissingConfigError for project_id, team", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".mate", "plugins.yaml")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("plugins.yaml written on refusal: %v", err)
	}
}

func TestAddBindingDisabledInstance(t *testing.T) {
	pkgs := packagesDir(t)
	root := workspace(t, "apiVersion: mate.plugins/v1\ninstances:\n  notes:\n    path: $PKG/skillsonly\n    enabled: false\n    projects: {blog: {}}\n")
	res, err := AddBinding(root, AddRequest{Instance: "notes", Path: filepath.Join(pkgs, "skillsonly"), Project: "shop"})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Disabled {
		t.Fatalf("result = %+v, want Disabled", res)
	}
	want := "apiVersion: mate.plugins/v1\ninstances:\n  notes:\n    path: " + filepath.Join(pkgs, "skillsonly") +
		"\n    enabled: false\n    projects:\n      blog: {}\n      shop: {}\n"
	if got := readPlugins(t, root); got != want {
		t.Fatalf("plugins.yaml =\n%s\nwant\n%s", got, want)
	}
}

func TestAddBindingNeedsAWorkspace(t *testing.T) {
	_, err := AddBinding(t.TempDir(), AddRequest{Instance: "x", Path: "/x", Project: "shop"})
	if err == nil || !strings.Contains(err.Error(), "is not a workspace") {
		t.Fatalf("err = %v", err)
	}
}
