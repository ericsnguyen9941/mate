package store

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
)

// methodGroup is one reason a method of the store changes
// (docs/plans/store-split-2026-11-06.md section 2): the files that hold the
// group, and the methods it held when the store was split by group
// (CON-239). Each file holds one group only, so two changes for different
// reasons land in different files.
//
// A new method needs no row: the file it is declared in names its group.
// A new file does: add it to the files of its group. A listed method that
// is renamed or removed must be renamed or removed here too.
type methodGroup struct {
	name    string
	files   []string
	methods []string
}

// methodGroups is the table TestStoreFilesKeepTheirGroup checks. The four
// path groups are G1 split by the directory the path is in. A group with no
// methods is a file that must hold none: workspace.go keeps only the handle
// and how it is opened.
var methodGroups = []methodGroup{
	{
		name:  "G1 paths of the workspace",
		files: []string{"layout.go"},
		methods: []string{
			"Root", "StateDir", "GCLog", "WorkspaceFile", "WorkspaceDoc", "WorkspaceCrewDoc",
			"PricingFile", "EnvFile", "JevLog", "ScreensDir", "WorktreesDir", "RepoDir",
			"ConsoleLockFile", "CoreLockFile", "CoreSocket", "CoreLogFile", "TUILockFile",
			"RecoverLockFile", "MigrateLockFile", "MigrateLog",
		},
	},
	{
		name:  "G1 paths of a project",
		files: []string{"layout_project.go"},
		methods: []string{
			"ProjectsDir", "ProjectDir", "ProjectHome", "ProjectFile", "ProjectDoc",
			"ProjectCrewDoc", "SentLog", "IncidentsLog",
		},
	},
	{
		name:  "G1 paths of a Mate",
		files: []string{"layout_mate.go"},
		methods: []string{
			"MateDir", "MateMeta", "AutoFlag", "ManualHold", "AutoCursorFile", "WakeCursorFile",
			"OutboxFile", "outboxLockFile", "MemoryFile", "BacklogFile", "CrewLearning", "LearningsDir",
		},
	},
	{
		name:  "G1 paths of a Crew",
		files: []string{"layout_crew.go"},
		methods: []string{
			"CrewsDir", "CrewMeta", "CrewStatus", "CrewStopLock", "CrewPRWatch", "CrewPRWatchLog",
			"CrewTeardownLog", "CrewDir", "CrewBrief", "CrewReport", "CrewHandback",
			"CrewReviewMeta", "CrewHarnessProfile", "CrewSessionInput", "WorktreeDir",
		},
	},
	{
		name:  "G2 workspace configuration",
		files: []string{"workspace_config.go", "env.go", "pricing.go"},
		methods: []string{
			"LaneCap", "snapshot", "Config", "Session", "RecordedRoot", "SetRoot", "SetSession",
			"Layout", "LayoutOld", "Defaults", "CrewAutocompact", "CrewContextHigh", "LoadConfig",
			"SaveConfig", "updateConfig", "LoadEnv", "LoadPricing", "SavePricing",
		},
	},
	{
		name:  "G3 projects and repos",
		files: []string{"project.go", "repo.go"},
		methods: []string{
			"Projects", "RequireProject", "Project", "AddProject", "EnsureProjectDoc",
			"RemoveProject", "LoadProject", "SaveProject", "normaliseProject", "RelRepo",
			"normaliseRepos", "checkReposUnclaimed", "AddRepo", "CheckAddRepo", "addRepoPlan",
			"RemoveRepo",
		},
	},
	{
		name:    "G4 Auto and Manual mode",
		files:   []string{"mode.go"},
		methods: []string{"Auto", "Held", "SetMode", "removeFlag", "SetAuto"},
	},
	{
		name:  "G5 Crew records",
		files: []string{"meta.go", "harness_profile.go", "crewdoc.go", "review.go", "scout.go"},
		methods: []string{
			"ReadCrewMeta", "WriteCrewMeta", "UpdateCrewMeta", "CrewIDs",
			"WriteCrewHarnessProfile", "ReadCrewHarnessProfile", "readHarnessProfile",
			"ReadCrewHarnessProfiles", "CrewRules", "ReplaceCrewBrief", "WriteCrewSessionInput",
			"AddCrewReview", "ReviewLinks", "CrewIsScout",
		},
	},
	{
		name:  "G6 Mate records",
		files: []string{"mate_meta.go", "mate_sessions.go", "backlog.go", "memory.go"},
		methods: []string{
			"ReadMateMeta", "WriteMateMeta", "MateHeadless", "ArchiveMateSession",
			"MateSessionArchives", "FreezeMateSession", "EditBacklog", "EditBacklogEntry",
			"ReplaceMemoryFile",
		},
	},
	{
		name:  "G7 append-only logs",
		files: []string{"log.go"},
		methods: []string{
			"AppendStatus", "AppendSent", "ReadStatus", "ReadSent", "SentLogSize",
			"AppendJevLog", "AppendIncident", "ReadIncidents",
		},
	},
	{
		name:  "G8 outbox and cursors",
		files: []string{"outbox.go", "auto.go"},
		methods: []string{
			"ReadOutbox", "UpdateOutbox", "writeOutbox", "ReadAutoCursor", "ReadWakeCursor",
			"WriteWakeCursor", "readCursor", "WriteAutoCursor", "writeCursor",
		},
	},
	{
		name: "G9 locks between processes",
		files: []string{
			"console_lock.go", "core_lock.go", "tui_lock.go", "crew_stop_lock.go",
			"recover_lock.go", "maintenance.go", "crew_send.go", "lock.go",
		},
		methods: []string{
			"LockConsole", "ConsoleHolder", "LockCore", "CoreHolder", "CoreRunning",
			"ConsoleRunning", "LockTUI", "TUIRunning", "forgetHarness", "LockCrewStop",
			"CrewStopping", "LockRecover", "TryMateMaintenance", "TryMateLock", "WithCrewSend",
		},
	},
	{
		name:  "G10 migrate from layout 1 to 2",
		files: []string{"migrate.go"},
		methods: []string{
			"LockMigrate", "MigrateLocked", "AppendMigrateLog", "ReadMigrateLog",
			"WriteCrewBrief", "SetLayoutProjectDirs",
		},
	},
	{
		name:    "G11 boundary and file writes",
		files:   []string{"boundary.go", "fsio.go"},
		methods: []string{"resolve", "Resolve", "mkdirAll", "writeFile", "appendLine", "seedOnce"},
	},
	{
		name:  "the handle: Workspace, Init, Open",
		files: []string{"workspace.go"},
	},
}

// groupedReceivers are the types whose methods the groups place.
var groupedReceivers = map[string]bool{"Workspace": true, "Paths": true, "Migration": true}

// layoutImports are all a path file may import: a path helper only joins
// strings, it never touches the disk.
var layoutImports = map[string]bool{"fmt": true, "path/filepath": true}

// storeMethod is one method of a grouped receiver as declared in a file.
type storeMethod struct {
	recv, name, file string
}

// scanStore parses the package's files and returns the methods of the
// grouped receivers and each file's imports.
func scanStore(files map[string]*ast.File) ([]storeMethod, map[string][]string) {
	var methods []storeMethod
	imports := map[string][]string{}
	for name, f := range files {
		for _, spec := range f.Imports {
			path, _ := strconv.Unquote(spec.Path.Value)
			imports[name] = append(imports[name], path)
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
			if id, ok := recv.(*ast.Ident); ok && groupedReceivers[id.Name] {
				methods = append(methods, storeMethod{recv: id.Name, name: fn.Name.Name, file: name})
			}
		}
	}
	sort.Slice(methods, func(i, j int) bool {
		if methods[i].file != methods[j].file {
			return methods[i].file < methods[j].file
		}
		return methods[i].name < methods[j].name
	})
	return methods, imports
}

// groupOfFileGroup is the group whose files include name.
func groupOfFileGroup(groups []methodGroup, name string) methodGroup {
	for _, g := range groups {
		if slices.Contains(g.files, name) {
			return g
		}
	}
	return methodGroup{}
}

func isLayoutFile(name string) bool {
	return strings.HasPrefix(name, "layout") && strings.HasSuffix(name, ".go")
}

// checkGroups returns what is wrong with files against groups, one line
// each, empty when every method is in a file of its group.
func checkGroups(groups []methodGroup, files map[string]*ast.File) []string {
	var problems []string
	groupOfFile := map[string]string{}
	groupOfMethod := map[string]methodGroup{}
	for _, g := range groups {
		for _, f := range g.files {
			if other, ok := groupOfFile[f]; ok {
				problems = append(problems, fmt.Sprintf("groups_test.go puts %s in both %q and %q: a file holds one group", f, other, g.name))
			}
			groupOfFile[f] = g.name
			if _, ok := files[f]; !ok {
				problems = append(problems, fmt.Sprintf("groups_test.go lists %s under %q, but the package has no such file: drop it from the table", f, g.name))
			}
		}
		for _, m := range g.methods {
			if other, ok := groupOfMethod[m]; ok {
				problems = append(problems, fmt.Sprintf("groups_test.go lists %s in both %q and %q", m, other.name, g.name))
			}
			groupOfMethod[m] = g
		}
	}

	methods, imports := scanStore(files)
	declared := map[string]bool{}
	for _, m := range methods {
		declared[m.name] = true
		group, ok := groupOfFile[m.file]
		if !ok {
			problems = append(problems, fmt.Sprintf("%s declares %s.%s but is in no group: add %s to the files of its group in methodGroups (groups_test.go), or move the method to a file of its group", m.file, m.recv, m.name, m.file))
			continue
		}
		if g, ok := groupOfMethod[m.name]; ok && !slices.Contains(g.files, m.file) {
			problems = append(problems, fmt.Sprintf("%s.%s is in %s (%s), but it belongs to %q: move it to %s", m.recv, m.name, m.file, group, g.name, strings.Join(g.files, " or ")))
		} else if len(groupOfFileGroup(groups, m.file).methods) == 0 {
			problems = append(problems, fmt.Sprintf("%s.%s is in %s, which holds no methods (%q): move it to a file of its group", m.recv, m.name, m.file, group))
		}
		if m.recv == "Paths" && !isLayoutFile(m.file) {
			problems = append(problems, fmt.Sprintf("Paths.%s is in %s: path helpers live in layout*.go", m.name, m.file))
		}
		if m.recv != "Paths" && isLayoutFile(m.file) {
			problems = append(problems, fmt.Sprintf("%s.%s is in %s: layout*.go holds only Paths methods, which do no I/O", m.recv, m.name, m.file))
		}
	}
	for _, g := range groups {
		for _, m := range g.methods {
			if !declared[m] {
				problems = append(problems, fmt.Sprintf("groups_test.go lists %s under %q, but no file declares it: rename or drop it in methodGroups", m, g.name))
			}
		}
	}
	for name, paths := range imports {
		if !isLayoutFile(name) {
			continue
		}
		for _, p := range paths {
			if !layoutImports[p] {
				problems = append(problems, fmt.Sprintf("%s imports %s: a path file imports only fmt and path/filepath", name, p))
			}
		}
	}
	sort.Strings(problems)
	return problems
}

// parseStoreFiles parses every non-test Go file of this package.
func parseStoreFiles(t *testing.T) map[string]*ast.File {
	t.Helper()
	names, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	files := map[string]*ast.File{}
	for _, name := range names {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, name, nil, parser.SkipObjectResolution)
		if err != nil {
			t.Fatal(err)
		}
		files[name] = f
	}
	return files
}

// Every method of Workspace, Paths and Migration is in a file of its group,
// so changes for different reasons do not meet in one file (CON-239).
func TestStoreFilesKeepTheirGroup(t *testing.T) {
	for _, p := range checkGroups(methodGroups, parseStoreFiles(t)) {
		t.Error(p)
	}
}

// The check itself: a method in the wrong file, a file in no group, a method
// in a file that holds none, a path helper outside layout*.go, a layout file
// that does I/O, and a stale row are each reported.
func TestCheckGroupsReportsMisplacedCode(t *testing.T) {
	groups := []methodGroup{
		{name: "paths", files: []string{"layout.go"}, methods: []string{"StateDir"}},
		{name: "logs", files: []string{"log.go"}, methods: []string{"AppendSent", "Gone"}},
		{name: "projects", files: []string{"project.go"}, methods: []string{"LoadProject"}},
		{name: "handle", files: []string{"workspace.go"}},
	}
	sources := map[string]string{
		"layout.go":    "package store\nimport \"os\"\nfunc (p Paths) StateDir() string { return \"\" }\nfunc (w *Workspace) Stat() { os.Stat(\"\") }\n",
		"log.go":       "package store\n",
		"project.go":   "package store\nfunc (w *Workspace) LoadProject() {}\nfunc (w *Workspace) AppendSent() {}\nfunc (p Paths) SentLog() string { return \"\" }\n",
		"tui.go":       "package store\nfunc (w *Workspace) LockTUI() {}\n",
		"workspace.go": "package store\nfunc Open() {}\nfunc (w *Workspace) Session() {}\n",
	}
	fset := token.NewFileSet()
	files := map[string]*ast.File{}
	for name, src := range sources {
		f, err := parser.ParseFile(fset, name, src, parser.SkipObjectResolution)
		if err != nil {
			t.Fatal(err)
		}
		files[name] = f
	}
	got := strings.Join(checkGroups(groups, files), "\n")
	for _, want := range []string{
		"Workspace.AppendSent is in project.go (projects), but it belongs to \"logs\": move it to log.go",
		"tui.go declares Workspace.LockTUI but is in no group",
		"Paths.SentLog is in project.go: path helpers live in layout*.go",
		"Workspace.Stat is in layout.go: layout*.go holds only Paths methods",
		"layout.go imports os",
		"groups_test.go lists Gone under \"logs\", but no file declares it",
		"Workspace.Session is in workspace.go, which holds no methods",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("checkGroups did not report %q; it reported:\n%s", want, got)
		}
	}
	if problems := checkGroups(groups[:1], map[string]*ast.File{"layout.go": files["layout.go"]}); len(problems) != 2 {
		t.Errorf("a layout file with one stray method and one stray import: got %d problems, want 2: %q", len(problems), problems)
	}
}
