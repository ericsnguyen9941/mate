#!/usr/bin/env python3
"""Run one preregistered pair per scenario with fresh Codex sessions.

Requires macOS sandbox-exec, Go, and an authenticated Codex CLI. No install,
merge or user-workspace mutation. Raw receipts remain in the requested output.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import signal
import shutil
import subprocess
import sys
import time
import uuid
from fixtures import CASES, prepare, oracle, run

REPO=Path(__file__).resolve().parents[2]
BASELINE='9e6a6836ac6c48ad1f551f45dd63810126a947c8'
ORDER=[('S1','A'),('S1','B'),('S2','B'),('S2','A'),('S3','A'),('S3','B'),('S4','B'),('S4','A'),('S5','A'),('S5','B'),('S6','B'),('S6','A')]


def digest(data):
    return hashlib.sha256(data.encode() if isinstance(data,str) else data).hexdigest()


def prompt(f, variant, renderer):
    root=Path(f['root']);scope='scoped' if f['prior'] else 'full'
    baseline=(REPO/'scripts/revieweval/testdata'/f'baseline-{scope}.md').read_text()
    replacements={'__WORKSPACE__':str(root),'__HEAD__':f['head'],'__BASE__':f['base'],'__PRIOR_HEAD__':f['prior'] or '', '__PRIOR_FINDING__':CASES[f['case']].get('prior_finding','')}
    for a,b in replacements.items(): baseline=baseline.replace(a,b)
    if variant=='A': body=baseline
    else:
        params={'Crew':'change','Head':f['head'],'Base':f['base'],'BriefPath':str(root/'.mate/projects/shop/crews/change/brief.md'),'HandbackPath':str(root/'.mate/projects/shop/crews/change/handback.md'),'BaseDescription':next(s for s in baseline.splitlines() if s.startswith('- Base '))}
        if f['prior']: params['Prior']={'ID':'r1','Head':f['prior'],'ReportPath':str(root/'.mate/projects/shop/crews/r1/report.md'),'Findings':[CASES[f['case']]['prior_finding']]}
        body=run([str(renderer)],root,input=json.dumps(params)).stdout
    common=f'''Act as the single independent reviewer. Your own checkout is {f['repo']}; it is already at the requested SHA. The task documents are under {root}/.mate. Return the complete Markdown review report as your final response. Temporary probes may be written in this checkout; preserve tracked product code. Work autonomously using shell tools. This review does not require status-file handoff. Use only this task's checkout and documents, with no other agents or reviewer sessions.\n\n'''
    return common+body


def main():
    parser=argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--out',type=Path,required=True)
    parser.add_argument('--model',default='gpt-6-luna')
    parser.add_argument('--effort',default='medium')
    parser.add_argument('--timeout',type=int,default=300)
    parser.add_argument('--prepare-only',action='store_true')
    args=parser.parse_args();out=args.out.resolve();out.mkdir(parents=True,exist_ok=False)
    renderer=out/'render';run(['go','build','-o',str(renderer),'./scripts/revieweval'],REPO)
    oracle_results={k:oracle(k,out) for k in CASES}
    (out/'oracles.json').write_text(json.dumps(oracle_results,indent=2))
    fixtures=[{**prepare(case),'variant':variant,'id':uuid.uuid4().hex[:12]} for case,variant in ORDER]
    (out/'manifest.json').write_text(json.dumps({'baseline':BASELINE,'model':args.model,'effort':args.effort,'timeout_seconds':args.timeout,'token_cap':None,'harness':run(['codex','--version'],REPO).stdout.strip(),'order':fixtures,'source_hashes':{str(p.relative_to(REPO)):digest(p.read_bytes()) for p in [REPO/'internal/mateassets/review.go',REPO/'scripts/revieweval/fixtures.py',REPO/'scripts/revieweval/run.py',*sorted((REPO/'scripts/revieweval/testdata').glob('*.md'))]}},indent=2))
    print('Prepared 12 independent runs; oracles validated.',flush=True)
    if args.prepare_only:return
    results=[]
    for f in fixtures:
        rid=f['id'];root=Path(f['root']);repo=Path(f['repo'])
        codex_home=root/'codex-home';codex_home.mkdir(mode=0o700)
        auth_source=Path(os.environ.get('CODEX_HOME',str(Path.home()/'.codex')))/'auth.json'
        shutil.copy2(auth_source,codex_home/'auth.json')
        (codex_home/'auth.json').chmod(0o600)
        scratch=root/'tmp';scratch.mkdir()
        text=prompt(f,f['variant'],renderer);(out/f'{rid}.prompt.md').write_text(text)
        # OS-enforced blindness: candidate cannot read the protocol, gold,
        # source of the treatment, receipts, or another candidate's checkout.
        denied={str(REPO),str(out),'/Users/erics/app/mate',*[x['root'] for x in fixtures if x['root']!=f['root']],*[x['root'] for x in oracle_results.values()]}
        profile='(version 1)\n(allow default)\n(deny file-write*)\n(allow file-write* (subpath '+json.dumps(str(root))+') (literal "/dev/null") (literal "/dev/tty"))\n'+''.join('(deny file-read-data (subpath '+json.dumps(p)+'))\n' for p in sorted(denied))
        profile_path=out/f'{rid}.sb';profile_path.write_text(profile)
        denied_read=run(['sandbox-exec','-f',str(profile_path),'/bin/cat',str(out/'manifest.json')],repo,check=False)
        if denied_read.returncode==0:raise RuntimeError('candidate can read the private manifest')
        response=repo/'response.md'
        cmd=['sandbox-exec','-f',str(profile_path),'codex','exec','--ignore-user-config','--ignore-rules','--json','--model',args.model,'-c',f'model_reasoning_effort="{args.effort}"','-c','approval_policy="never"','-c','project_doc_max_bytes=0','--sandbox','danger-full-access','--cd',str(repo),'--output-last-message',str(response),'-']
        env=os.environ.copy()
        for key in ('RECORDS_TOKEN','RECORDS_URL'):env.pop(key,None)
        env['GOCACHE']=str(repo/'.cache/go-build')
        env['CODEX_HOME']=str(codex_home)
        env['TMPDIR']=str(scratch)+'/'
        env['TMPPREFIX']=str(scratch/'zsh')
        start=time.monotonic();timed_out=False
        print(f'Start {rid} ({len(results)+1}/12)',flush=True)
        proc=subprocess.Popen(cmd,cwd=repo,env=env,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE,text=True,start_new_session=True)
        try:stdout,stderr=proc.communicate(text,timeout=args.timeout)
        except subprocess.TimeoutExpired:
            timed_out=True;os.killpg(proc.pid,signal.SIGTERM)
            try:stdout,stderr=proc.communicate(timeout=10)
            except subprocess.TimeoutExpired:
                os.killpg(proc.pid,signal.SIGKILL);stdout,stderr=proc.communicate()
        except BaseException:
            os.killpg(proc.pid,signal.SIGTERM)
            proc.wait(timeout=10)
            raise
        (out/f'{rid}.jsonl').write_text(stdout)
        (out/f'{rid}.stderr').write_text(stderr)
        rollout=''.join(p.read_text() for p in sorted((codex_home/'sessions').rglob('*.jsonl')))
        (out/f'{rid}.rollout.jsonl').write_text(rollout)
        (codex_home/'auth.json').unlink(missing_ok=True)
        elapsed=round(time.monotonic()-start,3)
        events=[]
        for line in (out/f'{rid}.jsonl').read_text().splitlines():
            try:events.append(json.loads(line))
            except json.JSONDecodeError:pass
        usage=[e.get('usage') for e in events if e.get('type')=='turn.completed']
        commands=[e['item'] for e in events if e.get('type')=='item.completed' and e.get('item',{}).get('type')=='command_execution']
        tool_events=[]
        for line in rollout.splitlines():
            try:e=json.loads(line)
            except json.JSONDecodeError:continue
            if e.get('payload',{}).get('type') in ('function_call','function_call_output','custom_tool_call','custom_tool_call_output'):tool_events.append(e)
        report=response.read_text() if response.exists() else ''
        (out/f'{rid}.report.md').write_text(report)
        tracked=run(['git','diff','HEAD'],repo).stdout
        head_after=run(['git','rev-parse','HEAD'],repo).stdout.strip()
        probes={}
        for p in repo.rglob('*'):
            if not p.is_file() or any(x in p.relative_to(repo).parts for x in ('.git','.cache')):continue
            if p.stat().st_size<=100000 and p.name!='response.md':
                try:probes[str(p.relative_to(repo))]=p.read_text()
                except UnicodeDecodeError:pass
        result={'id':rid,'elapsed_seconds':elapsed,'exit':proc.returncode,'timeout':timed_out,'prompt_sha256':digest(text),'usage':usage,'commands':commands,'tool_events':tool_events,'head_after':head_after,'tracked_diff':tracked,'files_after':probes}
        (out/f'{rid}.receipt.json').write_text(json.dumps(result,indent=2))
        results.append({k:v for k,v in result.items() if k not in ('commands','tool_events','files_after')})
        (out/'results.json').write_text(json.dumps(results,indent=2))
        print(f'Finish {rid}: {elapsed}s exit={proc.returncode} usage={usage}',flush=True)
    print(str(out),flush=True)

if __name__=='__main__':main()
