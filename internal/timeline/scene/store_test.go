package scene

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/nguyenngocanh94/mate/internal/db"
)

// Save writes only what changed: the projection reruns whole on every pass,
// and rewriting every stored row was most of the core's disk writes.
func TestSaveWritesOnlyTheRowsThatChanged(t *testing.T) {
	handle, err := db.OpenPath(filepath.Join(t.TempDir(), "mate.db"))
	if err != nil {
		t.Fatalf("OpenPath: %v", err)
	}
	t.Cleanup(func() { _ = handle.Close() })
	ctx := context.Background()
	for _, id := range []string{mateID, crewID} {
		if _, err := handle.SQL().Exec(`INSERT INTO actor(id, project, kind, name) VALUES (?, ?, 'crew', ?)`, id, project, id); err != nil {
			t.Fatalf("insert actor: %v", err)
		}
	}
	at := time.Date(2026, 11, 15, 9, 0, 0, 0, time.UTC)
	tr := func(id string, to State, detail string) Transition {
		return Transition{ID: project + "|" + id, Project: project, ActorID: crewID, To: to, At: at, EventID: 1, Detail: detail}
	}
	save := func(rows ...Transition) int64 {
		t.Helper()
		tx, err := handle.SQL().BeginTx(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		var before, after int64
		_ = tx.QueryRow(`SELECT total_changes()`).Scan(&before)
		if err := Save(ctx, tx, project, rows); err != nil {
			t.Fatalf("Save: %v", err)
		}
		_ = tx.QueryRow(`SELECT total_changes()`).Scan(&after)
		if err := tx.Commit(); err != nil {
			t.Fatal(err)
		}
		return after - before
	}

	a, b, c := tr("a", AtDeskWorking, ""), tr("b", Arriving, ""), tr("c", AtDeskWorking, "x")
	if n := save(a, b, c); n != 3 {
		t.Fatalf("first save changed %d rows, want 3", n)
	}
	if n := save(a, b, c); n != 0 {
		t.Fatalf("saving the same rows changed %d rows, want 0", n)
	}
	b2, d := tr("b", Arriving, "moved"), tr("d", Arriving, "")
	if n := save(a, b2, d); n != 3 {
		t.Fatalf("save with one changed, one new and one gone row changed %d rows, want 3", n)
	}

	rows, err := handle.SQL().Query(`SELECT id, detail, target_actor_id FROM transition ORDER BY id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var got []string
	for rows.Next() {
		var id, detail string
		var target sql.NullString
		if err := rows.Scan(&id, &detail, &target); err != nil {
			t.Fatal(err)
		}
		got = append(got, id+"="+detail)
		if target.Valid {
			t.Errorf("row %s has target %q, want NULL", id, target.String)
		}
	}
	want := []string{"shop|a=", "shop|b=moved", "shop|d="}
	if len(got) != len(want) {
		t.Fatalf("rows = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("rows = %v, want %v", got, want)
		}
	}
}
