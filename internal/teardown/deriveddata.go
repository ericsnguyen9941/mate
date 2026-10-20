package teardown

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DerivedData deletes the Xcode build folders made for the worktree. Xcode
// keys `~/Library/Developer/Xcode/DerivedData/<Name>-<hash>` on the
// project's path and records it as `WorkspacePath` in the folder's
// info.plist, so a worktree gets folders of its own that nothing ever
// reuses once it is gone. A repo with its own DerivedData root (plantrip's
// `-derivedDataPath`) cleans that in its teardown hook.
type DerivedData struct {
	// Root is the DerivedData directory. Empty or missing sweeps nothing.
	Root string
}

// Sweep implements Sweep.
func (d DerivedData) Sweep(_ context.Context, t Target) ([]string, error) {
	if d.Root == "" {
		return nil, nil
	}
	entries, err := os.ReadDir(d.Root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("DerivedData: %v", err)
	}
	var removed, failed []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(d.Root, e.Name())
		if !t.Owns(WorkspacePath(filepath.Join(dir, "info.plist"))) {
			continue
		}
		if err := os.RemoveAll(dir); err != nil {
			failed = append(failed, fmt.Sprintf("DerivedData %s: %v", e.Name(), err))
			continue
		}
		removed = append(removed, "DerivedData "+e.Name())
	}
	if len(failed) > 0 {
		return removed, fmt.Errorf("%s", strings.Join(failed, "; "))
	}
	return removed, nil
}

// WorkspacePath is the `WorkspacePath` string of a DerivedData info.plist,
// or "" when the file is missing or says none. Xcode writes it as an XML
// plist: a top-level dict of <key> and value elements.
func WorkspacePath(plist string) string {
	f, err := os.Open(plist)
	if err != nil {
		return ""
	}
	defer f.Close()
	dec := xml.NewDecoder(f)
	key, wantValue := "", false
	for {
		tok, err := dec.Token()
		if err != nil {
			return ""
		}
		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		var text string
		if start.Name.Local == "key" || start.Name.Local == "string" {
			if err := dec.DecodeElement(&text, &start); err != nil {
				return ""
			}
		}
		switch {
		case start.Name.Local == "key":
			key, wantValue = text, true
		case wantValue && key == "WorkspacePath":
			if start.Name.Local == "string" {
				return strings.TrimSpace(text)
			}
			return ""
		default:
			wantValue = false
		}
	}
}
