package pi

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/nguyenngocanh94/mate/internal/harness"
)

func TestPiMCPConfigAndGlobalCollision(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(AgentDirEnv, dir)
	servers := []harness.MCPServer{{Name: "work-probe", Type: "stdio", Command: []string{"./server", "quoted \"arg"}, Cwd: "/package dir"}, {Name: "remote", Type: "http", URL: "https://example.com/mcp"}}
	doc, err := MCPConfig(servers)
	if err != nil {
		t.Fatal(err)
	}
	var decoded struct {
		Servers map[string]struct {
			Command, Cwd, URL, Exposure string
			Args                        []string
		} `json:"mcpServers"`
	}
	if err := json.Unmarshal(doc, &decoded); err != nil {
		t.Fatal(err)
	}
	got := decoded.Servers["work-probe"]
	if got.Command != "./server" || got.Cwd != "/package dir" || got.Exposure != "direct" || !slices.Equal(got.Args, servers[0].Command[1:]) {
		t.Fatalf("%+v", got)
	}
	if err := os.WriteFile(filepath.Join(dir, "mcp.json"), []byte(`{"mcpServers":{"work_probe":{"url":"https://other.com"}}}`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := checkGlobalMCP(servers); err == nil {
		t.Fatal("global namespace collision accepted")
	}
	if _, err := MCPConfig(append(servers, harness.MCPServer{Name: "work_probe", Type: "http", URL: "https://example.com"})); err == nil {
		t.Fatal("normalized collision accepted")
	}
	if _, err := MCPConfig([]harness.MCPServer{{Name: "has.dot", Type: "http", URL: "https://example.com"}}); err == nil {
		t.Fatal("invalid native pi name accepted")
	}
}
func TestPiMateExtensionEvents(t *testing.T) {
	req := crewRequest(t)
	req.Role = harness.RoleMate
	logPath := filepath.Join(req.Cwd, "payloads.jsonl")
	req.Binary = filepath.Join(req.Cwd, "hook-test")
	shell := "#!/bin/sh\ncat >> " + harness.ShellQuote(logPath) + "\nprintf '\\n' >> " + harness.ShellQuote(logPath) + "\nprintf DIGEST\n"
	if err := os.WriteFile(req.Binary, []byte(shell), 0755); err != nil {
		t.Fatal(err)
	}
	prep, err := (Pi{}).Prepare(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range prep.Files {
		if err := os.MkdirAll(filepath.Dir(f.Path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(f.Path, f.Data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	// Replay real event shapes through the embedded module. The fake binary
	// records payloads, which catches duplicate receipts and source mapping.
	script := `
import {writeFileSync,readFileSync} from 'node:fs';
const configPath=process.argv[1];const extension=process.argv[2];
const events={};const servers=[];
const pi={registerMcpServer:(n,s)=>servers.push([n,s]),on:(n,f)=>{events[n]=f}};
const {default:install}=await import(extension);install(pi);
const ctx={cwd:process.cwd(),hasUI:false,sessionManager:{getSessionId:()=>"test-session",getSessionFile:()=>"test.jsonl"}};
await events.session_start({reason:'startup'},ctx);
await events.input({text:'⟦mate⟧ hello'},ctx);
const first=await events.before_agent_start({},ctx);
if(first?.message?.content!=="DIGEST"||first.message.customType!=='mate-recall')throw Error('no recall');
if(await events.before_agent_start({},ctx))throw Error('duplicate recall');
await events.session_compact({},ctx);
if(!(await events.before_agent_start({},ctx))?.message)throw Error('no compact recall');
await events.agent_end({messages:[{role:'assistant',stopReason:'stop',content:[{type:'text',text:'answer'}]}]},ctx);
await events.session_start({reason:'resume'},ctx);
if(!(await events.before_agent_start({},ctx))?.message)throw Error('no resume recall');
`
	cmd := exec.Command("node", "--input-type=module", "-e", script, launchConfigPath(req.Cwd), extensionPath(req.Cwd))
	cmd.Dir = req.Cwd
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("embedded extension events: %v %s", err, out)
	}
	b, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	var payloads []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		var p map[string]any
		if err := json.Unmarshal([]byte(line), &p); err != nil {
			t.Fatal(err)
		}
		payloads = append(payloads, p)
	}
	if len(payloads) != 5 || payloads[0]["prompt"] != "⟦mate⟧ hello" || payloads[1]["source"] != "startup" || payloads[2]["source"] != "compact" || payloads[3]["last_assistant_message"] != "answer" || payloads[4]["source"] != "resume" {
		t.Fatalf("payloads: %+v", payloads)
	}

}

func TestPiMateFilesProtectOwnershipAndLaunch(t *testing.T) {
	req := crewRequest(t)
	req.Role = harness.RoleMate
	req.Binary = "/bin/mate"
	files, err := mateFiles(req)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if err := os.MkdirAll(filepath.Dir(f.Path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(f.Path, f.Data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := checkMateFiles(req.Cwd, nil, false); err != nil {
		t.Fatal(err)
	}
	if err := checkMateFiles(req.Cwd, nil, true); err == nil {
		t.Fatal("accepted config for a different launch intent")
	}
	if err := os.WriteFile(extensionPath(req.Cwd), []byte("user extension"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := mateFiles(req); err == nil {
		t.Fatal("overwrote an unowned extension")
	}
	if err := checkMateFiles(req.Cwd, nil, false); err == nil {
		t.Fatal("accepted a changed extension")
	}
	if err := os.WriteFile(launchConfigPath(req.Cwd), []byte(`{"binary":"/missing"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if changed, err := (piHooks{}).Repoint("/new/mate", req.Cwd, func(string) bool { return false }); changed || err == nil {
		t.Fatal("repointed an unowned config")
	}
}
