package store

import (
	"os"
	"path/filepath"
)

// mkdirAll creates dir and its parents after checking the boundary.
func (w *Workspace) mkdirAll(dir string) error {
	if _, err := w.resolve(dir); err != nil {
		return err
	}
	return os.MkdirAll(dir, 0o755)
}

// writeFile writes path atomically: a temp file in the same directory, then a
// rename. The temp file is removed if anything fails, so a failed write leaves
// no debris next to the real file.
func (w *Workspace) writeFile(path string, data []byte, perm os.FileMode) error {
	if _, err := w.resolve(path); err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := w.mkdirAll(dir); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	cleanup := func() {
		tmp.Close()
		os.Remove(name)
	}
	if _, err := tmp.Write(data); err != nil {
		cleanup()
		return err
	}
	if err := tmp.Chmod(perm); err != nil {
		cleanup()
		return err
	}
	if err := tmp.Sync(); err != nil {
		cleanup()
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(name)
		return err
	}
	if err := os.Rename(name, path); err != nil {
		os.Remove(name)
		return err
	}
	return nil
}

// appendLine writes line plus a newline to path in one O_APPEND write, holding
// an exclusive lock for the write.
func (w *Workspace) appendLine(path, line string) error {
	if _, err := w.resolve(path); err != nil {
		return err
	}
	if err := w.mkdirAll(filepath.Dir(path)); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := lockFile(f, true); err != nil {
		return err
	}
	defer unlockFile(f)
	_, err = f.WriteString(line + "\n")
	return err
}

// seedOnce writes body to path unless path already exists. It never
// overwrites: from the moment it exists the file is the captain's.
func (w *Workspace) seedOnce(path, body string) error {
	if _, err := w.resolve(path); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	_, werr := f.WriteString(body)
	if cerr := f.Close(); werr == nil {
		werr = cerr
	}
	return werr
}
