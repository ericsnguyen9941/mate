package store

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// pathCases pins what every Paths method returns for one fixed root. A path
// in `.mate/` or `.worktrees/` is a file format other processes and the
// captain's tools rely on, so a change here is a behaviour change, not a
// refactor. The keys are the method names; TestPathsTableCoversEveryMethod
// keeps the table complete.
func pathCases(p Paths) map[string]struct{ got, want string } {
	const project, crew = "shop", "c1"
	return map[string]struct{ got, want string }{
		"Root":             {p.Root(), "/ws"},
		"StateDir":         {p.StateDir(), "/ws/.mate"},
		"GCLog":            {p.GCLog(), "/ws/.mate/gc.log"},
		"WorkspaceFile":    {p.WorkspaceFile(), "/ws/.mate/workspace.yaml"},
		"WorkspaceDoc":     {p.WorkspaceDoc(), "/ws/.mate/WORKSPACE.md"},
		"WorkspaceCrewDoc": {p.WorkspaceCrewDoc(), "/ws/.mate/CREW.md"},
		"PricingFile":      {p.PricingFile(), "/ws/.mate/pricing.yaml"},
		"EnvFile":          {p.EnvFile(), "/ws/.mate/.env"},
		"JevLog":           {p.JevLog(), "/ws/.mate/jev.log"},
		"ScreensDir":       {p.ScreensDir(), "/ws/.mate/screens"},
		"ConsoleLockFile":  {p.ConsoleLockFile(), "/ws/.mate/console.lock"},
		"CoreLockFile":     {p.CoreLockFile(), "/ws/.mate/core.lock"},
		"CoreSocket":       {p.CoreSocket(), "/ws/.mate/core.sock"},
		"CoreLogFile":      {p.CoreLogFile(), "/ws/.mate/core.log"},
		"RecoverLockFile":  {p.RecoverLockFile(), "/ws/.mate/recover.lock"},
		"MigrateLockFile":  {p.MigrateLockFile(), "/ws/.mate/migrate.lock"},
		"MigrateLog":       {p.MigrateLog(), "/ws/.mate/migrate.log"},

		"ProjectsDir":    {p.ProjectsDir(), "/ws/.mate/projects"},
		"ProjectDir":     {p.ProjectDir(project), "/ws/.mate/projects/shop"},
		"ProjectHome":    {p.ProjectHome(project), "/ws/shop"},
		"ProjectFile":    {p.ProjectFile(project), "/ws/.mate/projects/shop/project.yaml"},
		"ProjectDoc":     {p.ProjectDoc(project), "/ws/.mate/projects/shop/PROJECT.md"},
		"ProjectCrewDoc": {p.ProjectCrewDoc(project), "/ws/.mate/projects/shop/CREW.md"},
		"SentLog":        {p.SentLog(project), "/ws/.mate/projects/shop/sent.log"},
		"IncidentsLog":   {p.IncidentsLog(project), "/ws/.mate/projects/shop/incidents.log"},

		"MateDir":        {p.MateDir(project), "/ws/.mate/projects/shop/mate"},
		"MateMeta":       {p.MateMeta(project), "/ws/.mate/projects/shop/mate/mate.meta"},
		"AutoFlag":       {p.AutoFlag(project), "/ws/.mate/projects/shop/mate/.auto"},
		"ManualHold":     {p.ManualHold(project), "/ws/.mate/projects/shop/mate/.manual"},
		"AutoCursorFile": {p.AutoCursorFile(project), "/ws/.mate/projects/shop/mate/.auto-cursor"},
		"WakeCursorFile": {p.WakeCursorFile(project), "/ws/.mate/projects/shop/mate/.wake-cursor"},
		"MemoryFile":     {p.MemoryFile(project), "/ws/.mate/projects/shop/mate/memory.md"},
		"BacklogFile":    {p.BacklogFile(project), "/ws/.mate/projects/shop/mate/backlog.md"},
		"OutboxFile":     {p.OutboxFile(project), "/ws/.mate/projects/shop/mate/.outbox"},
		"outboxLockFile": {p.outboxLockFile(project), "/ws/.mate/projects/shop/mate/.outbox.lock"},
		"LearningsDir":   {p.LearningsDir(project), "/ws/.mate/projects/shop/mate/learnings"},
		"CrewLearning":   {p.CrewLearning(project, crew), "/ws/.mate/projects/shop/mate/learnings/c1.md"},

		"CrewsDir":           {p.CrewsDir(project), "/ws/.mate/projects/shop/crews"},
		"CrewMeta":           {p.CrewMeta(project, crew), "/ws/.mate/projects/shop/crews/c1.meta"},
		"CrewStatus":         {p.CrewStatus(project, crew), "/ws/.mate/projects/shop/crews/c1.status"},
		"CrewPRWatch":        {p.CrewPRWatch(project, crew), "/ws/.mate/projects/shop/crews/c1.prwatch"},
		"CrewStopLock":       {p.CrewStopLock(project, crew), "/ws/.mate/projects/shop/crews/c1.stop.lock"},
		"CrewDir":            {p.CrewDir(project, crew), "/ws/.mate/projects/shop/crews/c1"},
		"CrewPRWatchLog":     {p.CrewPRWatchLog(project, crew), "/ws/.mate/projects/shop/crews/c1/prwatch.log"},
		"CrewTeardownLog":    {p.CrewTeardownLog(project, crew), "/ws/.mate/projects/shop/crews/c1/teardown.log"},
		"CrewBrief":          {p.CrewBrief(project, crew), "/ws/.mate/projects/shop/crews/c1/brief.md"},
		"CrewReport":         {p.CrewReport(project, crew), "/ws/.mate/projects/shop/crews/c1/report.md"},
		"CrewHandback":       {p.CrewHandback(project, crew), "/ws/.mate/projects/shop/crews/c1/handback.md"},
		"CrewSessionInput":   {p.CrewSessionInput(project, crew, 3), "/ws/.mate/projects/shop/crews/c1/session-3.md"},
		"CrewHarnessProfile": {p.CrewHarnessProfile(project, crew), "/ws/.mate/projects/shop/crews/c1/harness-profile.json"},
		"CrewReviewMeta":     {p.CrewReviewMeta(project, crew), "/ws/.mate/projects/shop/crews/c1/review.meta"},

		"WorktreesDir": {p.WorktreesDir(), "/ws/.worktrees"},
		"WorktreeDir":  {p.WorktreeDir(project, crew), "/ws/.worktrees/shop-c1"},
		"RepoDir":      {p.RepoDir("shop/app"), "/ws/shop/app"},
	}
}

func TestPathsPinned(t *testing.T) {
	p := Paths{root: filepath.FromSlash("/ws")}
	for name, c := range pathCases(p) {
		if got := filepath.ToSlash(c.got); got != c.want {
			t.Errorf("%s = %q, want %q", name, got, c.want)
		}
	}
}

// A workspace's helpers are its Paths' helpers: embedding promotes them, so
// `ws.CrewMeta(p, c)` keeps compiling and keeps naming the same file.
func TestWorkspacePromotesPaths(t *testing.T) {
	w := &Workspace{Paths: Paths{root: filepath.FromSlash("/ws")}}
	if got, want := w.CrewMeta("shop", "c1"), w.Paths.CrewMeta("shop", "c1"); got != want {
		t.Fatalf("Workspace.CrewMeta = %q, Paths.CrewMeta = %q", got, want)
	}
	if got := filepath.ToSlash(w.Root()); got != "/ws" {
		t.Fatalf("Workspace.Root = %q, want /ws", got)
	}
}

// Every method declared on Paths in a non-test file must be pinned above, so
// a new path helper cannot land without the path it names written down.
func TestPathsTableCoversEveryMethod(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	var declared []string
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, name, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range f.Decls {
			fn, ok := d.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || len(fn.Recv.List) != 1 {
				continue
			}
			recv := fn.Recv.List[0].Type
			if star, ok := recv.(*ast.StarExpr); ok {
				recv = star.X
			}
			if id, ok := recv.(*ast.Ident); ok && id.Name == "Paths" {
				declared = append(declared, fn.Name.Name)
			}
		}
	}
	sort.Strings(declared)
	cases := pathCases(Paths{root: "/ws"})
	for _, name := range declared {
		if _, ok := cases[name]; !ok {
			t.Errorf("Paths.%s has no row in pathCases (paths_test.go): add the path it returns", name)
		}
	}
	if len(declared) != len(cases) {
		t.Errorf("Paths declares %d methods, pathCases pins %d", len(declared), len(cases))
	}
}
