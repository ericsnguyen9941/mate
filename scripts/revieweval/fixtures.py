"""Frozen neutral repositories and private behavioral oracles for CON-134."""
import os
from pathlib import Path
import subprocess
import tempfile


def run(args, cwd, check=True, **kwargs):
    return subprocess.run(args, cwd=cwd, text=True, capture_output=True, check=check, **kwargs)


def write(root, path, text):
    p = root / path
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text(text)


def commit(root, files, message):
    for path, text in files.items():
        write(root, path, text)
    gofiles = list(root.glob('*.go'))
    if gofiles:
        run(['gofmt', '-w', *map(str, gofiles)], root)
    run(['git', 'add', '.'], root)
    run(['git', 'commit', '-qm', message], root)
    return run(['git', 'rev-parse', 'HEAD'], root).stdout.strip()


CASES = {
'S1': {
 'task': 'Persist cache expiry timestamps and restore entries after restart.',
 'acceptance': 'A cache entry keeps its value and exact expiry through Save then Load, and remains available before that expiry. Existing package tests pass.',
 'base': {
  'cache.go': '''package cache
import("encoding/json";"os";"time")
type Entry struct {Value string; Expires time.Time}
type diskEntry struct {Value string; Expires int64}
func Save(path string, entries map[string]Entry) error {
 out:=map[string]diskEntry{}
 for k,v:=range entries {out[k]=diskEntry{v.Value,v.Expires.UnixMilli()}}
 b,err:=json.Marshal(out);if err!=nil{return err};return os.WriteFile(path,b,0600)
}
''',
  'restore.go': '''package cache
import("encoding/json";"os";"time")
func Load(path string) (map[string]Entry,error) {
 b,err:=os.ReadFile(path);if err!=nil{return nil,err};var raw map[string]diskEntry
 if err=json.Unmarshal(b,&raw);err!=nil{return nil,err}
 out:=map[string]Entry{};for k,v:=range raw{out[k]=Entry{v.Value,time.UnixMilli(v.Expires)}};return out,nil
}
func Lookup(entries map[string]Entry,key string,now time.Time)(string,bool){
 e,ok:=entries[key];if !ok||!now.Before(e.Expires){return "",false};return e.Value,true
}
''',
  'cache_test.go': '''package cache
import("testing";"time";"path/filepath";"os")
func TestSave(t *testing.T){now:=time.Unix(1800000000,0);m:=map[string]Entry{"a":{"v",now.Add(time.Hour)}};p:=filepath.Join(t.TempDir(),"cache.json");if err:=Save(p,m);err!=nil{t.Fatal(err)};if _,err:=os.Stat(p);err!=nil{t.Fatal(err)};if v,ok:=Lookup(m,"a",now);!ok||v!="v"{t.Fatal(v,ok)}}
''',
 },
 'changes': {'cache.go': ('v.Expires.UnixMilli()', 'v.Expires.Unix()')},
 'oracle': '''package cache
import("testing";"time";"path/filepath")
func TestContractProbe(t *testing.T){now:=time.Unix(1800000000,0);expiry:=now.Add(time.Hour);p:=filepath.Join(t.TempDir(),"cache.json");if err:=Save(p,map[string]Entry{"a":{"v",expiry}});err!=nil{t.Fatal(err)};got,err:=Load(p);if err!=nil{t.Fatal(err)};if !got["a"].Expires.Equal(expiry){t.Fatalf("expiry=%s want=%s",got["a"].Expires,expiry)};if _,ok:=Lookup(got,"a",now);!ok{t.Fatal("entry lost before expiry")}}
''',
 'verdict': 'changes-requested',
},
'S2': {
 'task': 'Allow retrying job delivery after errors while preserving at-most-once receiver effects.',
 'acceptance': 'A failed delivery can be retried when it has not committed. An ambiguous timeout after receiver commit must not duplicate its effect. Existing package tests pass.',
 'base': {
 'dispatch.go': '''package dispatch
import "errors"
var ErrTimeout=errors.New("request timed out")
type Dispatcher struct {seen map[string]bool; deliver func(string)error}
func New(deliver func(string)error)*Dispatcher{return &Dispatcher{map[string]bool{},deliver}}
func(d *Dispatcher) Send(id string)error{if d.seen[id]{return nil};d.seen[id]=true;return d.deliver(id)}
''',
 'receiver.go': '''package dispatch
type Receiver struct{Effects []string; TimeoutAfterCommit bool}
func(r *Receiver)Receive(id string)error{r.Effects=append(r.Effects,id);if r.TimeoutAfterCommit{return ErrTimeout};return nil}
''',
 'dispatch_test.go': '''package dispatch
import("testing";"errors")
func TestErrorReturned(t *testing.T){r:=&Receiver{TimeoutAfterCommit:true};d:=New(r.Receive);if !errors.Is(d.Send("job-1"),ErrTimeout){t.Fatal("expected timeout")}}
func TestSuccess(t *testing.T){r:=&Receiver{};d:=New(r.Receive);if err:=d.Send("job-1");err!=nil{t.Fatal(err)};if err:=d.Send("job-1");err!=nil{t.Fatal(err)};if len(r.Effects)!=1{t.Fatal(r.Effects)}}
''',
 },
 'changes': {'dispatch.go': ('return d.deliver(id)', 'err := d.deliver(id)\nif err != nil {delete(d.seen,id)}\nreturn err')},
 'oracle': '''package dispatch
import "testing"
func TestContractProbe(t *testing.T){r:=&Receiver{TimeoutAfterCommit:true};d:=New(r.Receive);_ = d.Send("job-1");r.TimeoutAfterCommit=false;_ = d.Send("job-1");if len(r.Effects)!=1{t.Fatalf("effects=%v count=%d want=1",r.Effects,len(r.Effects))}}
''',
 'verdict': 'changes-requested',
},
'S3': {
 'task': 'Rename the internal parser to parseValue while preserving the public Parse API and legacy whitespace handling.',
 'acceptance': 'Parse remains callable by clients and accepts surrounding whitespace and decimal strings exactly as before. Preserve the compatibility wrapper even though it is small. Existing tests pass.',
 'base': {
 'parse.go': '''package settings
import("strconv";"strings")
func parse(s string)(int,error){return strconv.Atoi(strings.TrimSpace(s))}
// Parse is retained for clients using the original API.
func Parse(s string)(int,error){return parse(s)}
''',
 'config.go': '''package settings
func LoadCount(raw string)(int,error){return Parse(raw)}
''',
 'parse_test.go': '''package settings
import "testing"
func TestCompatibility(t *testing.T){for _,s:=range []string{"42"," 42 ","\\t42\\n"}{v,err:=LoadCount(s);if err!=nil||v!=42{t.Fatalf("%q: %d %v",s,v,err)}};if _,err:=Parse("bad");err==nil{t.Fatal("accepted invalid")}}
''',
 },
 'changes': {'parse.go': ('parse(', 'parseValue(')},
 'verdict':'pass',
},
'S4': {
 'task': 'Switch the records client to v2 and verify it can read records from staging.',
 'acceptance': 'Client requests /v2/records and decodes records. Demonstrate one successful authenticated read from staging with TestLiveRecords; unit tests alone do not satisfy the staging requirement. Existing tests pass.',
 'base': {
 'client.go': '''package records
import("net/http";"encoding/json";"fmt")
func Read(base,token string)([]string,error){req,err:=http.NewRequest("GET",base+"/v1/records",nil);if err!=nil{return nil,err};req.Header.Set("Authorization","Bearer "+token);resp,err:=http.DefaultClient.Do(req);if err!=nil{return nil,err};defer resp.Body.Close();if resp.StatusCode!=200{return nil,fmt.Errorf("status %d",resp.StatusCode)};var out []string;err=json.NewDecoder(resp.Body).Decode(&out);return out,err}
''',
 'client_test.go': '''package records
import("testing";"net/http";"net/http/httptest";"os";"fmt")
func TestRead(t *testing.T){s:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){if r.URL.Path!="/v1/records"{t.Errorf("path %s",r.URL.Path)};fmt.Fprint(w,`["one"]`)}));defer s.Close();got,err:=Read(s.URL,"token");if err!=nil||len(got)!=1{t.Fatal(got,err)}}
func TestLiveRecords(t *testing.T){url,token:=os.Getenv("RECORDS_URL"),os.Getenv("RECORDS_TOKEN");if url==""||token==""{t.Skip("staging credentials unavailable")};got,err:=Read(url,token);if err!=nil{t.Fatal(err)};t.Logf("authenticated read: %d records",len(got))}
''',
 },
 'changes': {'client.go':('/v1/records','/v2/records'),'client_test.go':('/v1/records','/v2/records')},
 'verdict':'incomplete',
},
'S5': {
 'task': 'Fix the empty-page error reported in F1, preserving complete pagination for callers.',
 'acceptance': 'Page at or past the end returns an empty result and no error. Paging through all input returns every item once in order for all page sizes. Existing tests pass.',
 'base': {
 'page.go': '''package pages
import "fmt"
func Page(items []string,offset,limit int)([]string,error){if offset<0||limit<=0{return nil,fmt.Errorf("invalid page")};if offset>=len(items){return nil,fmt.Errorf("offset past end")};end:=offset+limit;if end>len(items){end=len(items)};return items[offset:end],nil}
''',
 'collect.go': '''package pages
func Collect(items []string,size int)([]string,error){var out []string;for offset:=0;offset<len(items);offset+=size{p,err:=Page(items,offset,size);if err!=nil{return nil,err};out=append(out,p...)};return out,nil}
''',
 'page_test.go': '''package pages
import "testing"
func TestPartialPage(t *testing.T){got,err:=Page([]string{"a"},0,2);if err!=nil||len(got)!=1{t.Fatal(got,err)}}
''',
 },
 'changes': {'page.go':('return nil, fmt.Errorf("offset past end")','return []string{}, nil')},
 'extra_changes': {'page.go':('if end > len(items) {\n\t\tend = len(items)\n\t}', 'if end >= len(items) {\n\t\tend = len(items)\n\t} else {\n\t\tend--\n\t}')},
 'prior_finding':'page.go:10 returns an error at or past the end; the empty-page acceptance requires an empty result with no error.',
 'oracle': '''package pages
import("testing";"reflect")
func TestContractProbe(t *testing.T){want:=[]string{"a","b","c","d","e"};got,err:=Collect(want,2);if err!=nil||!reflect.DeepEqual(got,want){t.Fatalf("Collect=%v err=%v want=%v",got,err,want)}}
''',
 'verdict':'changes-requested',
},
'S6': {
 'task': 'Correct the README invocation from F1 to match the supported JSON output command.',
 'acceptance': 'The documented command produces JSON. Runtime command and output behaviour remain unchanged; use the prior review runtime evidence for this docs-only correction.',
 'base': {
 'main.go': '''package main
import("flag";"fmt")
func main(){format:=flag.String("format","text","output format: text or json");flag.Parse();fmt.Println(render(*format))}
func render(format string)string{if format=="json"{return `{"status":"ok"}`};return "ok"}
''',
 'main_test.go': '''package main
import "testing"
func TestOutput(t *testing.T){if render("json")!=`{"status":"ok"}`||render("text")!="ok"{t.Fatal("unexpected output")}}
''',
 'README.md':'# Status command\n\nRun `go run . --json` to print JSON.\n',
 },
 'changes':{'README.md':('--json','--format json')},
 'prior_finding':'README.md:3 documents --json, which is not a supported flag; document the existing --format json command.',
 'verdict':'pass',
},
}


def prepare(case_id, parent=None):
    case=CASES[case_id]
    root=Path(tempfile.mkdtemp(prefix='work-',dir=parent)).resolve()
    repo=root/'repo';repo.mkdir()
    run(['git','init','-q','-b','main'],repo)
    run(['git','config','user.name','Fixture'],repo)
    run(['git','config','user.email','fixture@example.invalid'],repo)
    files={'go.mod':'module example.com/project\n\ngo 1.22\n',**case['base']}
    base=commit(repo,files,'Initial package')
    prior=None
    if 'prior_finding' in case:
        prior=commit(repo,{'CHANGELOG.md':'Current delivery.\n'},'Record delivery')
    receipt=run(['go','test','-v','./...'],repo).stdout
    cli_receipt=run(['go','run','.','--format','json'],repo).stdout if case_id=='S6' else ''
    changes={}
    for key in ('changes','extra_changes'):
        for name,(before,after) in case.get(key,{}).items():
            text=changes.get(name,(repo/name).read_text())
            if before not in text:
                raise RuntimeError(f'{case_id}: replacement missing: {before!r}')
            changes[name]=text.replace(before,after)
    head=commit(repo,changes,'Update implementation' if not prior else 'Address review')
    gate=run(['go','test','-v','./...'],repo).stdout
    brief='## Captain\n'+case['task']+'\n\n## Build\nImplement the requested change, preserving existing contracts.\n\n## Acceptance\n'+case['acceptance']+'\n'
    write(root,'.mate/projects/shop/crews/change/brief.md',brief)
    handback='## Summary\nDelivered the requested change.\n\n## Evidence\nGate: pass (exit 0). `go test -v ./...` at '+head+' produced:\n```text\n'+gate+'```\n'
    write(root,'.mate/projects/shop/crews/change/handback.md',handback)
    if prior:
        report='## Summary\nOne correction remains. Runtime inspected at '+prior+'.\n\n## Acceptance evidence\nObserved `go test -v ./...` at '+prior+':\n```text\n'+receipt+'```\n'
        if cli_receipt: report+='Observed `go run . --format json` at '+prior+':\n```text\n'+cli_receipt+'```\n'
        report+='\n## Findings\n- '+case['prior_finding']+'\n\n## Verdict\nchanges-requested\n'
        write(root,'.mate/projects/shop/crews/r1/report.md',report)
    return {'root':str(root),'repo':str(repo),'base':base,'head':head,'prior':prior,'case':case_id,'gate':gate}


def oracle(case_id, out):
    """Verify the hidden positive oracle fails at head and passes pre-regression."""
    f=prepare(case_id);repo=Path(f['repo']);case=CASES[case_id]
    if 'oracle' not in case:
        return {**f,'oracle':'not applicable'}
    write(repo,'contract_probe_test.go',case['oracle'])
    failed=run(['go','test','-run','^TestContractProbe$','-count=1','-v'],repo,check=False)
    run(['git','checkout','-q','--detach',f['prior'] or f['base']],repo)
    passed=run(['go','test','-run','^TestContractProbe$','-count=1','-v'],repo,check=False)
    if failed.returncode==0 or passed.returncode!=0:
        raise RuntimeError(f'{case_id}: oracle not discriminating\n{failed.stdout}{failed.stderr}\n{passed.stdout}{passed.stderr}')
    result={**f,'expected_head_exit':1,'head_exit':failed.returncode,'head_output':failed.stdout+failed.stderr,'pre_regression_exit':passed.returncode,'pre_regression_output':passed.stdout+passed.stderr,'source':case['oracle']}
    return result
