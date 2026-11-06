package archtest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

// storeReader is a package that reads the store through narrow interfaces
// of its own and writes nothing through it: the rows of the table in
// docs/plans/store-split-2026-11-06.md section 3 whose "Ghi qua store?"
// column says no (step S7, CON-244).
type storeReader struct {
	pkg string
	// ifaces are the interfaces the package takes instead of
	// *store.Workspace, in its workspace.go.
	ifaces []string
	// whole is true when the whole package only reads the store: then none
	// of its non-test files holds a *store.Workspace, outside the
	// compile-time assertion `var _ Workspace = (*store.Workspace)(nil)` and
	// the nil check of rule 4 (a func noWorkspace that returns a bool).
	// It is false for spawn, which writes; only its ListCrews reads, through
	// CrewReader.
	whole bool
}

// storeReaders are the packages TestStoreReadersDoNotWrite holds to the
// rule. A package moves here when the plan's table says it writes nothing
// through the store; a package that writes (timeline, watch, autopilot,
// outbox, prwatch, hook) declares its writes in its interface instead.
var storeReaders = []storeReader{
	{pkg: "internal/box", ifaces: []string{"Workspace"}, whole: true},
	{pkg: "internal/query", ifaces: []string{"BoxWorkspace", "Workspace"}, whole: true},
	{pkg: "internal/dashboard", ifaces: []string{"Workspace"}, whole: true},
	{pkg: "internal/usage", ifaces: []string{"Workspace"}, whole: true},
	{pkg: "internal/learnings", ifaces: []string{"Workspace"}, whole: true},
	{pkg: "internal/db", ifaces: []string{"Workspace"}, whole: true},
	{pkg: "internal/spawn", ifaces: []string{"CrewReader"}},
}

const storePkg = module + "internal/store"

// storeWriteRule is the store's naming rule for writes, read from
// internal/store/groups_test.go so that it is written down once: the write
// prefixes, and the methods that write without one.
type storeWriteRule struct {
	prefixes   []string
	exceptions map[string]bool
}

// writes reports whether a store method named name writes: its name is a
// write prefix followed by an upper-case letter or nothing (Settings is
// not), or the store lists it as a write without a prefix.
func (r storeWriteRule) writes(name string) bool {
	if r.exceptions[name] {
		return true
	}
	for _, p := range r.prefixes {
		rest, ok := strings.CutPrefix(name, p)
		if ok && (rest == "" || unicode.IsUpper([]rune(rest)[0])) {
			return true
		}
	}
	return false
}

// loadStoreWriteRule reads writePrefixes and writesWithoutPrefix from
// internal/store/groups_test.go, where TestStoreWritesCarryAWriteName holds
// the store's exported methods to them.
func loadStoreWriteRule(t *testing.T) storeWriteRule {
	t.Helper()
	path := filepath.Join(root, "internal", "store", "groups_test.go")
	f, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	rule := storeWriteRule{exceptions: map[string]bool{}}
	str := func(e ast.Expr) string {
		lit, ok := e.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			t.Fatalf("%s: %T in the write rule is not a string literal", path, e)
		}
		s, _ := strconv.Unquote(lit.Value)
		return s
	}
	for _, d := range f.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs := spec.(*ast.ValueSpec)
			if len(vs.Names) != 1 || len(vs.Values) != 1 {
				continue
			}
			lit, ok := vs.Values[0].(*ast.CompositeLit)
			if !ok {
				continue
			}
			switch vs.Names[0].Name {
			case "writePrefixes":
				for _, e := range lit.Elts {
					rule.prefixes = append(rule.prefixes, str(e))
				}
			case "writesWithoutPrefix":
				for _, e := range lit.Elts {
					rule.exceptions[str(e.(*ast.KeyValueExpr).Key)] = true
				}
			}
		}
	}
	if len(rule.prefixes) == 0 {
		t.Fatalf("%s declares no writePrefixes list of string literals; this test reads the store's write rule from it", path)
	}
	return rule
}

// listedPackage is what `go list -export -deps -json` says of a package.
type listedPackage struct {
	ImportPath string
	Dir        string
	GoFiles    []string
	Export     string
}

// listExports runs `go list -export -deps` over pkgs and returns every
// package it names, the dependencies included, with the file of its export
// data. The go command compiles what its build cache does not hold yet.
func listExports(t *testing.T, pkgs []string) map[string]listedPackage {
	t.Helper()
	args := []string{"list", "-export", "-deps", "-json=ImportPath,Dir,GoFiles,Export", "--"}
	for _, p := range pkgs {
		args = append(args, module+p)
	}
	cmd := exec.Command("go", args...)
	cmd.Dir = root
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go %s: %v\n%s", strings.Join(args, " "), err, stderr.Bytes())
	}
	listed := map[string]listedPackage{}
	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		var p listedPackage
		if err := dec.Decode(&p); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			t.Fatalf("decode go list output: %v", err)
		}
		listed[p.ImportPath] = p
	}
	return listed
}

// TestStoreReadersDoNotWrite holds the packages of storeReaders to what the
// plan says of them: they read the store and write nothing through it.
// The check is on types, so an embedded interface counts with every method
// it brings, however deep:
//
//   - no interface of ifaces has a method that writes by the store's rule
//     (groups_test.go: a write prefix, or listed in writesWithoutPrefix);
//   - in a whole package, no non-test file holds a *store.Workspace outside
//     the compile-time assertion and the nil check, so no write is reached
//     through the concrete handle (store.Open included), and no write
//     method of the store is called through an interface of its own.
func TestStoreReadersDoNotWrite(t *testing.T) {
	rule := loadStoreWriteRule(t)
	var pkgs []string
	for _, r := range storeReaders {
		pkgs = append(pkgs, r.pkg)
	}
	listed := listExports(t, append(pkgs, "internal/store"))
	fset := token.NewFileSet()
	imp := importer.ForCompiler(fset, "gc", func(path string) (io.ReadCloser, error) {
		p, ok := listed[path]
		if !ok || p.Export == "" {
			return nil, fmt.Errorf("go list gave no export data for %s", path)
		}
		return os.Open(p.Export)
	})
	store, err := imp.Import(storePkg)
	if err != nil {
		t.Fatalf("import %s: %v", storePkg, err)
	}
	wsObj, ok := store.Scope().Lookup("Workspace").(*types.TypeName)
	if !ok {
		t.Fatalf("%s declares no type Workspace", storePkg)
	}
	storeMethods := types.NewMethodSet(types.NewPointer(wsObj.Type()))
	storeWrites := map[string]bool{}
	for i := range storeMethods.Len() {
		if name := storeMethods.At(i).Obj().Name(); rule.writes(name) {
			storeWrites[name] = true
		}
	}
	if len(storeWrites) == 0 {
		t.Fatalf("no method of *store.Workspace is a write by the store's rule; the rule was not read")
	}

	for _, r := range storeReaders {
		p, ok := listed[module+r.pkg]
		if !ok {
			t.Errorf("storeReaders lists %s, but go list does not know it: drop it here and from docs/architecture.md section 1", r.pkg)
			continue
		}
		pkg, files, info := typeCheck(t, fset, imp, p)
		for _, name := range r.ifaces {
			for _, problem := range writesInInterface(fset, pkg, name, rule) {
				t.Errorf("%s: %s", r.pkg, problem)
			}
		}
		if r.whole {
			for _, problem := range concreteStoreUses(fset, files, info, storeWrites) {
				t.Errorf("%s", problem)
			}
		}
	}
}

// typeCheck type-checks the non-test files of p from source, against the
// export data of its imports.
func typeCheck(t *testing.T, fset *token.FileSet, imp types.Importer, p listedPackage) (*types.Package, []*ast.File, *types.Info) {
	t.Helper()
	var files []*ast.File
	for _, name := range p.GoFiles {
		f, err := parser.ParseFile(fset, filepath.Join(p.Dir, name), nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatalf("parse %s: %v", name, err)
		}
		files = append(files, f)
	}
	info := &types.Info{
		Types:      map[ast.Expr]types.TypeAndValue{},
		Defs:       map[*ast.Ident]types.Object{},
		Selections: map[*ast.SelectorExpr]*types.Selection{},
	}
	pkg, err := (&types.Config{Importer: imp}).Check(p.ImportPath, fset, files, info)
	if err != nil {
		t.Fatalf("type-check %s: %v", p.ImportPath, err)
	}
	return pkg, files, info
}

// writesInInterface returns a line for each method of the interface name in
// pkg, its embedded interfaces' included, that writes by rule.
func writesInInterface(fset *token.FileSet, pkg *types.Package, name string, rule storeWriteRule) []string {
	obj, ok := pkg.Scope().Lookup(name).(*types.TypeName)
	if !ok {
		return []string{fmt.Sprintf("storeReaders names interface %s, but the package declares no type %s: fix the row in internal/archtest/store_readers_test.go", name, name)}
	}
	iface, ok := obj.Type().Underlying().(*types.Interface)
	if !ok {
		return []string{fmt.Sprintf("%s is not an interface", name)}
	}
	var problems []string
	for i := range iface.NumMethods() {
		m := iface.Method(i)
		if !rule.writes(m.Name()) {
			continue
		}
		problems = append(problems, fmt.Sprintf("%s.%s has %s, a write to the store (declared at %s). "+
			"The plan (docs/plans/store-split-2026-11-06.md section 3) has this package read the store and write nothing through it. "+
			"Move the write to a package that writes, or, if this package must write now, change the plan and drop it from storeReaders in internal/archtest/store_readers_test.go. "+
			"A method writes when its name starts with a prefix of writePrefixes in internal/store/groups_test.go followed by an upper-case letter, or when writesWithoutPrefix lists it",
			obj.Pkg().Name(), name, m.Name(), relPos(fset, m.Pos())))
	}
	return problems
}

// concreteStoreUses returns a line for each place in files that holds a
// *store.Workspace outside the compile-time assertion and the nil check,
// and for each call of a store write method through an interface.
func concreteStoreUses(fset *token.FileSet, files []*ast.File, info *types.Info, storeWrites map[string]bool) []string {
	var allowed []ast.Node
	for _, f := range files {
		for _, d := range f.Decls {
			switch d := d.(type) {
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					if vs, ok := spec.(*ast.ValueSpec); ok && d.Tok == token.VAR && allBlank(vs.Names) {
						allowed = append(allowed, vs)
					}
				}
			case *ast.FuncDecl:
				if isNilCheck(d) {
					allowed = append(allowed, d)
				}
			}
		}
	}
	inAllowed := func(pos token.Pos) bool {
		return slices.ContainsFunc(allowed, func(n ast.Node) bool { return n.Pos() <= pos && pos < n.End() })
	}

	seen := map[string]bool{}
	var problems []string
	report := func(pos token.Pos, msg string) {
		line := relPos(fset, pos) + ": " + msg
		if !seen[line] {
			seen[line] = true
			problems = append(problems, line)
		}
	}
	const holds = "holds a *store.Workspace. A package that only reads the store reaches it through its own Workspace interface, which shows it writes nothing; " +
		"only `var _ Workspace = (*store.Workspace)(nil)` and the nil check noWorkspace may name the concrete type (plan section 3, rules 1 and 4)"
	for expr, tv := range info.Types {
		if isStoreWorkspace(tv.Type) && !inAllowed(expr.Pos()) {
			report(expr.Pos(), holds)
		}
	}
	for id, obj := range info.Defs {
		if obj != nil && isStoreWorkspace(obj.Type()) && !inAllowed(id.Pos()) {
			report(id.Pos(), holds)
		}
	}
	for sel, s := range info.Selections {
		recv := s.Recv()
		switch {
		case isStoreWorkspace(recv):
			report(sel.Pos(), fmt.Sprintf("calls %s on a *store.Workspace; a reader of the store calls it through its Workspace interface", s.Obj().Name()))
		case types.IsInterface(recv) && storeWrites[s.Obj().Name()]:
			report(sel.Pos(), fmt.Sprintf("calls %s, a write of the store, through an interface: a package that only reads the store calls no write of it", s.Obj().Name()))
		}
	}
	sort.Strings(problems)
	return problems
}

// isStoreWorkspace reports whether t is store.Workspace or a pointer to it.
func isStoreWorkspace(t types.Type) bool {
	if p, ok := t.(*types.Pointer); ok {
		t = p.Elem()
	}
	n, ok := t.(*types.Named)
	return ok && n.Obj().Pkg() != nil && n.Obj().Pkg().Path() == storePkg && n.Obj().Name() == "Workspace"
}

func allBlank(names []*ast.Ident) bool {
	for _, n := range names {
		if n.Name != "_" {
			return false
		}
	}
	return len(names) > 0
}

// isNilCheck reports whether fn is the nil check of plan section 3, rule 4:
// a func noWorkspace that returns one bool, so the *store.Workspace it
// asserts does not leave it.
func isNilCheck(fn *ast.FuncDecl) bool {
	if fn.Recv != nil || fn.Name.Name != "noWorkspace" || fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		return false
	}
	res := fn.Type.Results.List[0]
	id, ok := res.Type.(*ast.Ident)
	return ok && id.Name == "bool" && len(res.Names) <= 1
}

// relPos is pos as file:line, the file relative to the repository root
// when it is under it.
func relPos(fset *token.FileSet, pos token.Pos) string {
	p := fset.Position(pos)
	if abs, err := filepath.Abs(root); err == nil {
		if rel, err := filepath.Rel(abs, p.Filename); err == nil && !strings.HasPrefix(rel, "..") {
			p.Filename = filepath.ToSlash(rel)
		}
	}
	return fmt.Sprintf("%s:%d", p.Filename, p.Line)
}

// The rule is read from the store's test file: a write prefix counts only
// before an upper-case letter or the end, and writesWithoutPrefix counts.
func TestStoreWriteRuleIsReadFromTheStore(t *testing.T) {
	rule := loadStoreWriteRule(t)
	for name, want := range map[string]bool{
		"SetMode": true, "AppendIncident": true, "LockCore": true, "TUIRunning": true,
		"Settings": false, "ReadStatus": false, "LoadConfig": false, "CrewStopLock": false,
	} {
		if got := rule.writes(name); got != want {
			t.Errorf("writes(%q) = %v, want %v", name, got, want)
		}
	}
}
