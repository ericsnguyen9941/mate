package pi

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

//go:embed mate-extension.mjs
var mateExtension []byte

func extensionPath(cwd string) string    { return filepath.Join(cwd, ".pi", "mate-extension.mjs") }
func launchConfigPath(cwd string) string { return filepath.Join(cwd, ".pi", "mate-launch.json") }

type piHooks struct{}

func (piHooks) Own(string, string) []harness.OwnHook { return nil }
func (piHooks) ReviewOwn(_ context.Context, _ harness.HookReviewPane, screen string, _ []harness.OwnHook) error {
	return &harness.ScreenRefusal{Screen: screen, Reason: "pi does not ask for hook trust; refusing an unexpected hook dialog"}
}
func (piHooks) DigestMaxBytes() int   { return 48000 }
func (piHooks) BareSessionHook() bool { return false }
func (piHooks) Repoint(binary, cwd string, exists func(string) bool) (bool, error) {
	path := launchConfigPath(cwd)
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	var doc map[string]any
	if err = json.Unmarshal(b, &doc); err != nil {
		return false, err
	}
	if doc["mateOwned"] != float64(1) {
		return false, fmt.Errorf("pi Mate will not repoint unowned file %s", path)
	}
	old, _ := doc["binary"].(string)
	if old == "" || exists(old) {
		return false, nil
	}
	doc["binary"] = binary
	b, err = json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return false, err
	}
	if err := replaceLaunchConfig(path, append(b, '\n')); err != nil {
		return false, err
	}
	return true, nil
}

// replaceLaunchConfig keeps the previous config intact until the whole new
// document is written. The temporary file stays beside it for atomic rename.
func replaceLaunchConfig(path string, data []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".mate-launch-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}
func mateFiles(req harness.PrepareRequest) ([]harness.LaunchFile, error) {
	if req.Binary == "" || !filepath.IsAbs(req.Binary) {
		return nil, fmt.Errorf("pi Mate requires an absolute mate binary")
	}
	for _, path := range []string{extensionPath(req.Cwd), launchConfigPath(req.Cwd)} {
		old, err := os.ReadFile(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		owned := bytes.HasPrefix(old, []byte("// Mate-owned pi"))
		if path == launchConfigPath(req.Cwd) {
			var doc struct {
				Owned int `json:"mateOwned"`
			}
			owned = json.Unmarshal(old, &doc) == nil && doc.Owned == 1
		}
		if !owned {
			return nil, fmt.Errorf("pi Mate would overwrite unowned file %s", path)
		}
	}
	b, err := MCPConfig(req.MCPServers)
	if err != nil {
		return nil, err
	}
	var doc map[string]any
	if err := json.Unmarshal(b, &doc); err != nil {
		return nil, fmt.Errorf("pi Mate rendered MCP config: %w", err)
	}
	doc["mateOwned"] = 1
	doc["binary"] = req.Binary
	doc["resume"] = req.ResumeSessionID != ""
	b, err = json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	return []harness.LaunchFile{{Path: extensionPath(req.Cwd), Data: mateExtension}, {Path: launchConfigPath(req.Cwd), Data: append(b, '\n')}}, nil
}

func checkMateFiles(cwd string, servers []harness.MCPServer, resume bool) error {
	ext, err := os.ReadFile(extensionPath(cwd))
	if err != nil {
		return err
	}
	if !bytes.Equal(ext, mateExtension) {
		return fmt.Errorf("pi Mate extension is missing or changed")
	}
	b, err := os.ReadFile(launchConfigPath(cwd))
	if err != nil {
		return err
	}
	var doc struct {
		Binary  string          `json:"binary"`
		Resume  bool            `json:"resume"`
		Owned   int             `json:"mateOwned"`
		Servers json.RawMessage `json:"mcpServers"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		return err
	}
	expected, err := MCPConfig(servers)
	if err != nil {
		return err
	}
	var want map[string]json.RawMessage
	if err := json.Unmarshal(expected, &want); err != nil {
		return fmt.Errorf("pi Mate rendered MCP config: %w", err)
	}
	var actualMap, expectedMap any
	if err := json.Unmarshal(doc.Servers, &actualMap); err != nil {
		return fmt.Errorf("pi Mate launch mcpServers: %w", err)
	}
	if err := json.Unmarshal(want["mcpServers"], &expectedMap); err != nil {
		return fmt.Errorf("pi Mate expected mcpServers: %w", err)
	}
	actual, err := json.Marshal(actualMap)
	if err != nil {
		return err
	}
	target, err := json.Marshal(expectedMap)
	if err != nil {
		return err
	}
	if doc.Owned != 1 || !filepath.IsAbs(doc.Binary) || doc.Resume != resume || !bytes.Equal(actual, target) {
		return fmt.Errorf("pi Mate launch config does not match this launch")
	}
	return nil
}
