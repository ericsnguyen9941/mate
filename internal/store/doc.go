// Package store is the only package that reads and writes `.mate/` and
// `.worktrees/`. Every other package - project CRUD, spawn, box, watch,
// console - goes through it, so the on-disk layout of docs/mvp.md section 3 is
// described in exactly one place.
//
// Ownership rule. The console process is the only writer of state: it creates
// the workspace, registers projects, writes `workspace.yaml`, `project.yaml`,
// the `.meta` files and the `.auto` flag. Agents do not write state. The one
// exception is `crews/<id>.status`, which a crew appends to with a plain
// `echo "state: one line" >> $MATE_STATUS`; that is why the status log is
// append-only, line oriented and locked rather than rewritten. `sent.log` is
// append-only for the same reason: the console and the Mate hooks both add
// lines to it and nobody edits it. Anything that is not an append is an atomic
// temp-file-plus-rename write, so a reader never sees half a file.
//
// Path boundary. Every path this package writes is resolved through
// filepath.EvalSymlinks on its deepest existing ancestor and must land inside
// the workspace root. A symlink under `.mate/` that points outside the
// workspace makes the write fail with *BoundaryError and nothing is created.
// Registered repositories must be inside the workspace too; they need not be
// direct children of it.
//
// Files. Each file holds the methods of one responsibility group, so changes
// made for different reasons land in different files
// (docs/plans/store-split-2026-11-06.md section 2). groups_test.go keeps it
// that way; a new file is added to its table.
//
//   - Paths, string joins only: layout.go (the workspace's own files),
//     layout_project.go, layout_mate.go, layout_crew.go.
//   - workspace.yaml, .env and pricing.yaml: workspace_config.go, env.go,
//     pricing.go.
//   - Projects and repos: project.go, repo.go; project.yaml's schema is in
//     project_config.go.
//   - Auto and Manual mode: mode.go.
//   - Crew records: meta.go, harness_profile.go, crewdoc.go, review.go,
//     scout.go.
//   - Mate records: mate_meta.go, mate_sessions.go, backlog.go, memory.go.
//   - Append-only logs (`.status`, `sent.log`, `incidents.log`, `jev.log`):
//     log.go.
//   - Outbox and cursors: outbox.go, auto.go.
//   - Locks between processes: console_lock.go, core_lock.go, tui_lock.go,
//     crew_stop_lock.go, recover_lock.go, maintenance.go, crew_send.go,
//     lock.go.
//   - Migrating layout 1 to 2: migrate.go (the Migration view, which only
//     `mate migrate` uses).
//   - The boundary and the file writers: boundary.go, fsio.go.
//
// workspace.go holds the handle and how it is opened (Init, Open,
// OpenForMigrate), and no method.
package store
