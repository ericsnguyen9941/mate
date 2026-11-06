package usage

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const module = "github.com/nguyenngocanh94/mate/"

// TestUsageDoesNotReachTheDashboard walks this package's imports through
// every package of the module they lead to, the way the linker would, and
// fails on internal/dashboard or a package under it: the dashboard is to
// reuse the ledger by importing usage, so usage must not link it (CON-234;
// `--why`, which reads through dashboard.Server, is in internal/usage/why).
// Test files are left out.
func TestUsageDoesNotReachTheDashboard(t *testing.T) {
	root := filepath.Join("..", "..")
	via := map[string]string{"internal/usage": ""}
	queue := []string{"internal/usage"}
	for len(queue) > 0 {
		pkg := queue[0]
		queue = queue[1:]
		for _, imp := range packageImports(t, filepath.Join(root, filepath.FromSlash(pkg))) {
			rel, ok := strings.CutPrefix(imp, module)
			if !ok {
				continue
			}
			if rel == "internal/dashboard" || strings.HasPrefix(rel, "internal/dashboard/") {
				t.Errorf("internal/usage reaches %s: %s -> %s", imp, chain(via, pkg), rel)
				continue
			}
			if _, seen := via[rel]; !seen {
				via[rel] = pkg
				queue = append(queue, rel)
			}
		}
	}
	// usage reaches spawn, store, db and timeline/read at least; fewer
	// means the walk read nothing.
	if len(via) < 5 {
		t.Fatalf("walked %d packages; the import walk found nothing to check", len(via))
	}
}

// packageImports is every import of the package in dir, over all of its
// non-test files whatever their build constraints: a file built only on
// another system links on that system.
func packageImports(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	var out []string
	fset := token.NewFileSet()
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, filepath.Join(dir, name), nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("parse %s: %v", filepath.Join(dir, name), err)
		}
		for _, imp := range file.Imports {
			path, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				t.Fatalf("%s: unquote import %s: %v", name, imp.Path.Value, err)
			}
			out = append(out, path)
		}
	}
	return out
}

// chain is the import path from internal/usage to pkg, for the failure
// message: which package let the dashboard in.
func chain(via map[string]string, pkg string) string {
	path := []string{pkg}
	for from := via[pkg]; from != ""; from = via[from] {
		path = append([]string{from}, path...)
	}
	return strings.Join(path, " -> ")
}
