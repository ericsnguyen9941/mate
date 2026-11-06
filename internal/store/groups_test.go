package store

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

// methodGroup is one reason a method of the store changes
// (docs/plans/store-split-2026-11-06.md section 2): the files that hold the
// group, and the methods it held when the store was split by group
// (CON-239). Each file holds one group only, so two changes for different
// reasons land in different files.
//
// A new method needs no row: the file it is declared in names its group.
// A new file does: add it to the files of its group. A listed method that
// is renamed or removed must be renamed or removed here too. kind says
// what the group's methods may do to the disk (CON-244); a group that holds
// no methods has none.
type methodGroup struct {
	name    string
	kind    groupKind
	files   []string
	methods []string
}

// methodGroups is the table TestStoreFilesKeepTheirGroup checks, one method
// per line so that a rename touches one line. The four path groups are G1
// split by what the path belongs to: the workspace, a project, its Mate or
// one of its Crews. A group with no methods is a file that must hold none:
// workspace.go keeps only the handle and how it is opened.
var methodGroups = []methodGroup{
	{
		name:  "G1 paths of the workspace",
		kind:  pathKind,
		files: []string{"layout.go"},
		methods: []string{
			"Root",
			"StateDir",
			"GCLog",
			"WorkspaceFile",
			"WorkspaceDoc",
			"WorkspaceCrewDoc",
			"PricingFile",
			"EnvFile",
			"JevLog",
			"ScreensDir",
			"WorktreesDir",
			"RepoDir",
			"ConsoleLockFile",
			"CoreLockFile",
			"CoreSocket",
			"CoreLogFile",
			"TUILockFile",
			"RecoverLockFile",
			"MigrateLockFile",
			"MigrateLog",
		},
	},
	{
		name:  "G1 paths of a project",
		kind:  pathKind,
		files: []string{"layout_project.go"},
		methods: []string{
			"ProjectsDir",
			"ProjectDir",
			"ProjectHome",
			"ProjectFile",
			"ProjectDoc",
			"ProjectCrewDoc",
			"SentLog",
			"IncidentsLog",
		},
	},
	{
		name:  "G1 paths of a Mate",
		kind:  pathKind,
		files: []string{"layout_mate.go"},
		methods: []string{
			"MateDir",
			"MateMeta",
			"AutoFlag",
			"ManualHold",
			"AutoCursorFile",
			"WakeCursorFile",
			"OutboxFile",
			"outboxLockFile",
			"MemoryFile",
			"BacklogFile",
			"CrewLearning",
			"LearningsDir",
		},
	},
	{
		name:  "G1 paths of a Crew",
		kind:  pathKind,
		files: []string{"layout_crew.go"},
		methods: []string{
			"CrewsDir",
			"CrewMeta",
			"CrewStatus",
			"CrewStopLock",
			"CrewPRWatch",
			"CrewPRWatchLog",
			"CrewTeardownLog",
			"CrewDir",
			"CrewBrief",
			"CrewReport",
			"CrewHandback",
			"CrewReviewMeta",
			"CrewHarnessProfile",
			"CrewSessionInput",
			"WorktreeDir",
		},
	},
	{
		name:  "G2 workspace configuration",
		kind:  recordKind,
		files: []string{"workspace_config.go", "env.go", "pricing.go"},
		methods: []string{
			"LaneCap",
			"snapshot",
			"Config",
			"Session",
			"RecordedRoot",
			"SetRoot",
			"SetSession",
			"Layout",
			"LayoutOld",
			"Defaults",
			"CrewAutocompact",
			"CrewContextHigh",
			"LoadConfig",
			"SaveConfig",
			"updateConfig",
			"LoadEnv",
			"LoadPricing",
			"SavePricing",
		},
	},
	{
		name:  "G3 projects and repos",
		kind:  recordKind,
		files: []string{"project.go", "repo.go"},
		methods: []string{
			"Projects",
			"RequireProject",
			"Project",
			"AddProject",
			"EnsureProjectDoc",
			"RemoveProject",
			"LoadProject",
			"SaveProject",
			"normaliseProject",
			"RelRepo",
			"normaliseRepos",
			"checkReposUnclaimed",
			"AddRepo",
			"CheckAddRepo",
			"addRepoPlan",
			"RemoveRepo",
		},
	},
	{
		name:  "G4 Auto and Manual mode",
		kind:  recordKind,
		files: []string{"mode.go"},
		methods: []string{
			"Auto",
			"Held",
			"SetMode",
			"removeFlag",
			"SetAuto",
		},
	},
	{
		name:  "G5 Crew records",
		kind:  recordKind,
		files: []string{"meta.go", "harness_profile.go", "crewdoc.go", "review.go", "scout.go"},
		methods: []string{
			"ReadCrewMeta",
			"WriteCrewMeta",
			"UpdateCrewMeta",
			"CrewIDs",
			"WriteCrewHarnessProfile",
			"ReadCrewHarnessProfile",
			"readHarnessProfile",
			"ReadCrewHarnessProfiles",
			"CrewRules",
			"ReplaceCrewBrief",
			"WriteCrewSessionInput",
			"AddCrewReview",
			"ReviewLinks",
			"CrewIsScout",
		},
	},
	{
		name:  "G6 Mate records",
		kind:  recordKind,
		files: []string{"mate_meta.go", "mate_sessions.go", "backlog.go", "memory.go"},
		methods: []string{
			"ReadMateMeta",
			"WriteMateMeta",
			"MateHeadless",
			"ArchiveMateSession",
			"MateSessionArchives",
			"FreezeMateSession",
			"EditBacklog",
			"EditBacklogEntry",
			"ReplaceMemoryFile",
		},
	},
	{
		name:  "G7 append-only logs",
		kind:  recordKind,
		files: []string{"log.go"},
		methods: []string{
			"AppendStatus",
			"AppendSent",
			"ReadStatus",
			"ReadSent",
			"SentLogSize",
			"AppendJevLog",
			"AppendIncident",
			"ReadIncidents",
		},
	},
	{
		name:  "G8 outbox and cursors",
		kind:  recordKind,
		files: []string{"outbox.go", "auto.go"},
		methods: []string{
			"ReadOutbox",
			"UpdateOutbox",
			"writeOutbox",
			"ReadAutoCursor",
			"ReadWakeCursor",
			"WriteWakeCursor",
			"readCursor",
			"WriteAutoCursor",
			"writeCursor",
		},
	},
	{
		name: "G9 locks between processes",
		kind: recordKind,
		files: []string{
			"console_lock.go", "core_lock.go", "tui_lock.go", "crew_stop_lock.go",
			"recover_lock.go", "maintenance.go", "crew_send.go", "lock.go",
		},
		methods: []string{
			"LockConsole",
			"ConsoleHolder",
			"LockCore",
			"CoreHolder",
			"CoreRunning",
			"ConsoleRunning",
			"LockTUI",
			"TUIRunning",
			"forgetHarness",
			"LockCrewStop",
			"CrewStopping",
			"LockRecover",
			"TryMateMaintenance",
			"TryMateLock",
			"WithCrewSend",
		},
	},
	{
		name:  "G10 migrate from layout 1 to 2",
		kind:  recordKind,
		files: []string{"migrate.go"},
		methods: []string{
			"LockMigrate",
			"MigrateLocked",
			"AppendMigrateLog",
			"ReadMigrateLog",
			"WriteCrewBrief",
			"SetLayoutProjectDirs",
		},
	},
	{
		name:  "G11 boundary and file writes",
		kind:  recordKind,
		files: []string{"boundary.go", "fsio.go"},
		methods: []string{
			"resolve",
			"Resolve",
			"mkdirAll",
			"writeFile",
			"appendLine",
			"seedOnce",
		},
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

func isLayoutFile(name string) bool {
	return strings.HasPrefix(name, "layout") && strings.HasSuffix(name, ".go")
}

// checkGroups returns what is wrong with files against groups, one line
// each, empty when every method is in a file of its group.
func checkGroups(groups []methodGroup, files map[string]*ast.File) []string {
	var problems []string
	groupOfFile := map[string]methodGroup{}
	groupOfMethod := map[string]methodGroup{}
	for _, g := range groups {
		for _, f := range g.files {
			if other, ok := groupOfFile[f]; ok {
				problems = append(problems, fmt.Sprintf("groups_test.go puts %s in both %q and %q: a file holds one group", f, other.name, g.name))
			}
			groupOfFile[f] = g
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
			problems = append(problems, fmt.Sprintf("%s.%s is in %s (%s), but it belongs to %q: move it to %s", m.recv, m.name, m.file, group.name, g.name, strings.Join(g.files, " or ")))
		} else if len(group.methods) == 0 {
			problems = append(problems, fmt.Sprintf("%s.%s is in %s, which holds no methods (%q): move it to a file of its group", m.recv, m.name, m.file, group.name))
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
// that does I/O, a stale row, a file in two groups and a listed file that
// does not exist are each reported.
func TestCheckGroupsReportsMisplacedCode(t *testing.T) {
	groups := []methodGroup{
		{name: "paths", files: []string{"layout.go"}, methods: []string{"StateDir"}},
		{name: "logs", files: []string{"log.go"}, methods: []string{"AppendSent", "Gone"}},
		{name: "projects", files: []string{"project.go"}, methods: []string{"LoadProject"}},
		{name: "handle", files: []string{"workspace.go"}},
		{name: "twice", files: []string{"log.go", "gone.go"}},
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
		"groups_test.go puts log.go in both \"logs\" and \"twice\"",
		"groups_test.go lists gone.go under \"twice\", but the package has no such file",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("checkGroups did not report %q; it reported:\n%s", want, got)
		}
	}
	if problems := checkGroups(groups[:1], map[string]*ast.File{"layout.go": files["layout.go"]}); len(problems) != 2 {
		t.Errorf("a layout file with one stray method and one stray import: got %d problems, want 2: %q", len(problems), problems)
	}
}

// groupKind is what the methods of a group may do to the disk: the kind
// column of docs/plans/store-split-2026-11-06.md, step S7.
type groupKind int

const (
	// pathKind methods only join strings: none writes, and none carries a
	// write prefix.
	pathKind groupKind = iota + 1
	// recordKind methods read and write records. The naming rule tells the
	// two apart: an exported method that writes carries a write prefix, one
	// that only reads does not. No group today only reads.
	recordKind
)

// writePrefixes are the prefixes of an exported method that writes the
// disk (plan step S7). A name is a write name when it is one of these
// followed by an upper-case letter or by nothing, so Settings or Address
// would not be. A package that only reads the store declares no write name
// in its Workspace interface: internal/archtest/store_readers_test.go
// checks that, and reads this list and writesWithoutPrefix from this file.
var writePrefixes = []string{
	"Add", "Append", "Archive", "Edit", "Ensure", "Freeze", "Lock", "Remove",
	"Replace", "Save", "Set", "Try", "Update", "With", "Write",
}

// writesWithoutPrefix are the exported methods that write although their
// name carries no write prefix, each with why. internal/archtest counts
// every entry as a write in a reader's interface. Prefer a write name to a
// new entry.
var writesWithoutPrefix = map[string]string{
	"TUIRunning": "an orphan TUI's harness that has ended is forgotten: tui.lock's record is rewritten without its harness_pid (CON-227)",
}

// isWriteName reports whether name starts with a write prefix followed by
// an upper-case letter or nothing.
func isWriteName(name string) bool {
	for _, p := range writePrefixes {
		rest, ok := strings.CutPrefix(name, p)
		if ok && (rest == "" || unicode.IsUpper([]rune(rest)[0])) {
			return true
		}
	}
	return false
}

// osWrites are the functions of package os that change the disk. os.OpenFile
// changes it when its flags ask to write (osWriteFlags); a lock file opened
// O_CREATE|O_RDWR only to be locked is not a write.
var osWrites = map[string]bool{
	"Chmod": true, "Chtimes": true, "Create": true, "CreateTemp": true,
	"Link": true, "Mkdir": true, "MkdirAll": true, "MkdirTemp": true,
	"Remove": true, "RemoveAll": true, "Rename": true, "Symlink": true,
	"Truncate": true, "WriteFile": true, "Chown": true, "Lchown": true,
}

var osWriteFlags = map[string]bool{"O_WRONLY": true, "O_APPEND": true, "O_TRUNC": true}

// fileWrites are the methods of an *os.File that write it. WriteAt and
// Truncate count on any receiver (a lock's holder record is rewritten
// through a field, l.f). The others count on a variable the same function
// set from os.OpenFile, os.Create or os.CreateTemp, so a strings.Builder's
// WriteString does not; so does passing such a variable first to another
// package's writer (fmt.Fprintf, json.NewEncoder, io.Copy...).
var fileWrites = map[string]bool{"Write": true, "WriteAt": true, "WriteString": true, "Truncate": true, "ReadFrom": true}

// fileOpeners are the functions of package os whose first result is a file
// that fileWrites can write.
var fileOpeners = map[string]bool{"OpenFile": true, "Create": true, "CreateTemp": true}

// writerFuncs are the functions of other packages that write the file passed
// as their first argument, matched by name.
var writerFuncs = map[string]bool{
	"Fprint": true, "Fprintf": true, "Fprintln": true, "NewEncoder": true,
	"NewWriter": true, "Copy": true, "CopyN": true, "WriteString": true,
}

// diskWrites returns, for each function and method of the package that
// writes the disk, directly or through the functions it calls, one chain of
// calls from it down to the write ("TUIRunning → forgetHarness → Write →
// WriteAt"). Calls are matched by name only, whatever the receiver, so it
// errs towards writing: a name that one type writes under counts as a write
// for every caller.
func diskWrites(files map[string]*ast.File) map[string]string {
	type body struct {
		node    *ast.BlockStmt
		imports map[string]bool // the names the file imports packages under
		osName  string
	}
	bodies := map[string][]body{}
	for _, name := range slices.Sorted(maps.Keys(files)) {
		f := files[name]
		imports := map[string]bool{}
		osName := ""
		for _, spec := range f.Imports {
			path, _ := strconv.Unquote(spec.Path.Value)
			name := path[strings.LastIndex(path, "/")+1:]
			if spec.Name != nil {
				name = spec.Name.Name
			}
			imports[name] = true
			if path == "os" {
				osName = name
			}
		}
		for _, d := range f.Decls {
			if fn, ok := d.(*ast.FuncDecl); ok && fn.Body != nil {
				bodies[fn.Name.Name] = append(bodies[fn.Name.Name], body{fn.Body, imports, osName})
			}
		}
	}

	writes := map[string]string{}
	calls := map[string][]string{}
	for _, name := range slices.Sorted(maps.Keys(bodies)) {
		for _, b := range bodies[name] {
			fileVars := openedFiles(b.node, b.osName)
			isFileVar := func(e ast.Expr) bool {
				id, ok := e.(*ast.Ident)
				return ok && fileVars[id.Name]
			}
			write := func(what string) {
				if _, ok := writes[name]; !ok {
					writes[name] = what
				}
			}
			ast.Inspect(b.node, func(n ast.Node) bool {
				switch n := n.(type) {
				case *ast.CallExpr:
					sel, ok := n.Fun.(*ast.SelectorExpr)
					if !ok {
						break
					}
					x, _ := sel.X.(*ast.Ident)
					switch {
					case x != nil && b.osName != "" && x.Name == b.osName:
						if osWrites[sel.Sel.Name] || (sel.Sel.Name == "OpenFile" && len(n.Args) > 1 && opensToWrite(n.Args[1], b.osName)) {
							write("os." + sel.Sel.Name)
						}
					case x != nil && b.imports[x.Name]:
						if writerFuncs[sel.Sel.Name] && len(n.Args) > 0 && isFileVar(n.Args[0]) {
							write(x.Name + "." + sel.Sel.Name)
						}
					case sel.Sel.Name == "WriteAt" || sel.Sel.Name == "Truncate":
						write(sel.Sel.Name)
					case fileWrites[sel.Sel.Name] && isFileVar(sel.X):
						write(sel.Sel.Name)
					}
				case *ast.SelectorExpr:
					if x, ok := n.X.(*ast.Ident); ok && b.imports[x.Name] {
						return false
					}
					if _, ok := bodies[n.Sel.Name]; ok {
						calls[name] = append(calls[name], n.Sel.Name)
					}
				case *ast.Ident:
					if _, ok := bodies[n.Name]; ok {
						calls[name] = append(calls[name], n.Name)
					}
				}
				return true
			})
		}
	}
	for changed := true; changed; {
		changed = false
		for _, name := range slices.Sorted(maps.Keys(calls)) {
			if _, ok := writes[name]; ok {
				continue
			}
			for _, callee := range calls[name] {
				if chain, ok := writes[callee]; ok && callee != name {
					writes[name] = callee + " → " + chain
					changed = true
					break
				}
			}
		}
	}
	for name, chain := range writes {
		writes[name] = name + " → " + chain
	}
	return writes
}

// openedFiles returns the names of the variables body sets from an
// os.OpenFile, os.Create or os.CreateTemp call.
func openedFiles(body ast.Node, osName string) map[string]bool {
	vars := map[string]bool{}
	if osName == "" {
		return vars
	}
	ast.Inspect(body, func(n ast.Node) bool {
		as, ok := n.(*ast.AssignStmt)
		if !ok || len(as.Rhs) != 1 || len(as.Lhs) == 0 {
			return true
		}
		call, ok := as.Rhs[0].(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if x, ok := sel.X.(*ast.Ident); ok && x.Name == osName && fileOpeners[sel.Sel.Name] {
			if id, ok := as.Lhs[0].(*ast.Ident); ok && id.Name != "_" {
				vars[id.Name] = true
			}
		}
		return true
	})
	return vars
}

// opensToWrite reports whether the flags of an os.OpenFile call name a flag
// that writes, or name no os.O_* flag at all (flags held in a variable):
// those count as a write.
func opensToWrite(flags ast.Expr, osName string) bool {
	named, writes := false, false
	ast.Inspect(flags, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			if x, ok := sel.X.(*ast.Ident); ok && x.Name == osName && strings.HasPrefix(sel.Sel.Name, "O_") {
				named = true
				writes = writes || osWriteFlags[sel.Sel.Name]
			}
		}
		return true
	})
	return writes || !named
}

// checkWriteNames returns what breaks the naming rule, one line each: an
// exported method of Workspace, Paths or Migration that writes the disk
// with no write prefix and no entry in exceptions, a method of a path group
// with a write name or a write, a group with methods and no kind, and an
// exception that no longer applies.
func checkWriteNames(groups []methodGroup, files map[string]*ast.File, exceptions map[string]string) []string {
	var problems []string
	kindOfFile := map[string]groupKind{}
	groupOfFile := map[string]string{}
	for _, g := range groups {
		for _, f := range g.files {
			kindOfFile[f] = g.kind
			groupOfFile[f] = g.name
		}
	}
	writes := diskWrites(files)
	methods, _ := scanStore(files)
	declared := map[string]bool{}
	noKind := map[string]bool{}
	for _, m := range methods {
		if g, ok := groupOfFile[m.file]; ok && kindOfFile[m.file] == 0 && !noKind[g] {
			noKind[g] = true
			problems = append(problems, fmt.Sprintf("group %q in methodGroups holds methods (%s) but has no kind: give it pathKind or recordKind", g, m.file))
		}
		if !ast.IsExported(m.name) {
			continue
		}
		declared[m.name] = true
		chain, writing := writes[m.name]
		named := isWriteName(m.name)
		_, excepted := exceptions[m.name]
		switch {
		case kindOfFile[m.file] == pathKind && (named || writing):
			problems = append(problems, fmt.Sprintf("%s.%s in %s is in a path group but has a write name or writes (%s): a path helper only joins strings", m.recv, m.name, m.file, chain))
		case writing && !named && !excepted:
			problems = append(problems, fmt.Sprintf("%s.%s writes the disk (%s) but its name has no write prefix (one of %s, then an upper-case letter): rename it, or add it to writesWithoutPrefix in groups_test.go with why, which counts it as a write in every reader's interface", m.recv, m.name, chain, strings.Join(writePrefixes, ", ")))
		case excepted && (named || !writing):
			problems = append(problems, fmt.Sprintf("writesWithoutPrefix lists %s, but it has a write prefix or no longer writes: drop it", m.name))
		}
	}
	for _, name := range slices.Sorted(maps.Keys(exceptions)) {
		if !declared[name] {
			problems = append(problems, fmt.Sprintf("writesWithoutPrefix lists %s, but no file declares it as an exported method: drop it", name))
		}
	}
	sort.Strings(problems)
	return problems
}

// An exported method that writes the disk says so in its name, so a
// reader's Workspace interface shows whether it writes (CON-244, plan step
// S7). LoadConfig only reads: it re-reads workspace.yaml into memory.
func TestStoreWritesCarryAWriteName(t *testing.T) {
	for _, p := range checkWriteNames(methodGroups, parseStoreFiles(t), writesWithoutPrefix) {
		t.Error(p)
	}
}

// The naming check itself: a read name that writes through a helper, a
// method value, an open file (WriteAt, WriteString, fmt.Fprintln) or
// OpenFile flags it cannot read, a path with a write name, a group with no
// kind and stale exceptions are each reported; a lock file opened to be
// locked, a builder's WriteString and a write name are not.
func TestCheckWriteNamesReportsHiddenWrites(t *testing.T) {
	groups := []methodGroup{
		{name: "paths", kind: pathKind, files: []string{"layout.go"}, methods: []string{"StateDir"}},
		{name: "logs", kind: recordKind, files: []string{"log.go"}, methods: []string{"ReadStatus"}},
		{name: "no kind", files: []string{"mode.go"}, methods: []string{"Auto"}},
	}
	sources := map[string]string{
		"layout.go": "package store\nfunc (p Paths) StateDir() string { return \"\" }\nfunc (p Paths) SetDir() string { return \"\" }\n",
		"log.go": `package store
import (
	"fmt"
	"os"
	"strings"
)
func (w *Workspace) ReadStatus() { os.OpenFile("", os.O_CREATE|os.O_RDWR, 0); var b strings.Builder; b.WriteString("") }
func (w *Workspace) ReadSent() { w.store("") }
func (w *Workspace) Cursor() { f, _ := os.OpenFile("", os.O_RDWR, 0); f.WriteAt(nil, 0) }
func (w *Workspace) Tail() { do(w.flush) }
func (w *Workspace) AppendSent() { w.store("") }
func (w *Workspace) Count() {}
func (w *Workspace) Holder() { f, _ := os.OpenFile("", os.O_RDWR, 0); f.WriteString("") }
func (w *Workspace) Note() { f, _ := os.OpenFile("", os.O_RDWR, 0); fmt.Fprintln(f, "") }
func (w *Workspace) Flags(flag int) { os.OpenFile("", flag, 0) }
func (w *Workspace) store(s string) { os.OpenFile(s, os.O_WRONLY|os.O_APPEND, 0) }
func (w *Workspace) flush() { os.Rename("", "") }
func do(f func()) { f() }
`,
		"mode.go": "package store\nfunc (w *Workspace) Auto() bool { return false }\n",
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
	exceptions := map[string]string{"AppendSent": "stale", "Count": "stale", "Gone": "stale"}
	problems := checkWriteNames(groups, files, exceptions)
	got := strings.Join(problems, "\n")
	for _, want := range []string{
		"Workspace.ReadSent writes the disk (ReadSent → store → os.OpenFile) but its name has no write prefix",
		"Workspace.Cursor writes the disk (Cursor → WriteAt)",
		"Workspace.Tail writes the disk (Tail → flush → os.Rename)",
		"Workspace.Holder writes the disk (Holder → WriteString)",
		"Workspace.Note writes the disk (Note → fmt.Fprintln)",
		"Workspace.Flags writes the disk (Flags → os.OpenFile)",
		"Paths.SetDir in layout.go is in a path group but has a write name",
		"group \"no kind\" in methodGroups holds methods (mode.go) but has no kind",
		"writesWithoutPrefix lists AppendSent, but it has a write prefix or no longer writes",
		"writesWithoutPrefix lists Count, but it has a write prefix or no longer writes",
		"writesWithoutPrefix lists Gone, but no file declares it",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("checkWriteNames did not report %q; it reported:\n%s", want, got)
		}
	}
	if len(problems) != 11 {
		t.Errorf("got %d problems, want 11 (ReadStatus, Auto and StateDir are fine):\n%s", len(problems), got)
	}
}

func TestIsWriteName(t *testing.T) {
	for name, want := range map[string]bool{
		"SetMode": true, "Set": true, "WriteCrewMeta": true, "LockCore": true, "TryMateLock": true,
		"Settings": false, "Address": false, "Lockout": false, "ReadStatus": false, "LoadConfig": false,
		"Session": false, "Withdraw": false, "CrewStopLock": false,
	} {
		if got := isWriteName(name); got != want {
			t.Errorf("isWriteName(%q) = %v, want %v", name, got, want)
		}
	}
}
