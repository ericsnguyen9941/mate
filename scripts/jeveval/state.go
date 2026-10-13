package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/nguyenngocanh94/mate/internal/harness/catalog"
	"github.com/nguyenngocanh94/mate/internal/notice"
	"github.com/nguyenngocanh94/mate/internal/screen"
	"github.com/nguyenngocanh94/mate/internal/screen/fixture"
	"github.com/nguyenngocanh94/mate/internal/screen/jev"
)

// stateThreshold is the confidence the settle would ask of Jev.
const stateThreshold = 0.85

// StateRow is one screen's result on the pane-state question.
type StateRow struct {
	Screen
	// Label is the state the fixture observer reads; unknown means the
	// fixture abstains and the screen is not scored.
	Label     screen.PaneState
	Hash      string
	LatencyMS int64
	Err       string
	Answer    jev.Answer
}

func (r StateRow) Scored() bool { return r.Label != screen.PaneUnknown }

func (r StateRow) Match() bool {
	return r.Err == "" && r.Scored() && screen.PaneState(r.Answer.Choice) == r.Label
}

// stateCorpus is the startup captures, the injection screens and the
// composer captures, each with the state fixture.ObserveState gives it. The
// composer captures joined on 2026-10-13, when the observer began asking the
// question of a running pane it is about to call stale.
func stateCorpus(root string) ([]StateRow, error) {
	all, err := loadCorpus(root)
	if err != nil {
		return nil, err
	}
	reg := catalog.Default()
	var rows []StateRow
	for _, s := range all {
		if s.Source != sourceStartup && s.Source != sourceInjection && s.Source != sourceSend {
			continue
		}
		p, err := reg.Lookup(s.Kind)
		if err != nil {
			return nil, err
		}
		reading, err := fixture.New().(screen.StateObserver).ObserveState(context.Background(), p.Screen(), s.Text)
		if err != nil {
			return nil, err
		}
		rows = append(rows, StateRow{Screen: s, Label: reading.State})
	}
	return rows, nil
}

// runState is run for -question state: the same corpus discipline, one
// request per screen carrying only the state question.
func runState(w io.Writer, cassette string, replay, list bool) error {
	root, err := moduleRoot()
	if err != nil {
		return err
	}
	rows, err := stateCorpus(root)
	if err != nil {
		return err
	}
	if list {
		fmt.Fprintln(w, "| File | Source | Harness | Fixture state |")
		fmt.Fprintln(w, "| --- | --- | --- | --- |")
		for _, r := range rows {
			fmt.Fprintf(w, "| %s | %s | %s | %s |\n", r.File, r.Source, r.Kind, r.Label)
		}
		return nil
	}
	if replay && cassette == "" {
		return errors.New("-replay needs -cassette")
	}
	var key string
	var recorded map[string]manifestEntry
	if replay {
		if recorded, err = readManifest(cassette); err != nil {
			return err
		}
	} else if key, err = readKey(); err != nil {
		return err
	}
	client := &http.Client{Timeout: notice.Timeout,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}

	m := manifest{Model: notice.Model, Date: time.Now().Format("2006-01-02"), Prompt: jev.StatePromptHash()}
	var done []StateRow
	var notRecorded []string
	for _, row := range rows {
		body, err := jev.StateRequestBody(jev.Prepare(row.Text, key))
		if err != nil {
			return err
		}
		row.Hash = jev.HashOf(body)
		var data []byte
		if replay {
			e, ok := recorded[row.Hash]
			if !ok {
				notRecorded = append(notRecorded, row.File)
				continue
			}
			delete(recorded, row.Hash)
			row.LatencyMS, row.Err = e.LatencyMS, e.Error
			if e.Error == "" {
				if data, err = os.ReadFile(filepath.Join(cassette, row.Hash+".json")); err != nil {
					return err
				}
			}
		} else {
			ctx, cancel := context.WithTimeout(context.Background(), notice.Timeout)
			start := time.Now()
			data, err = call(ctx, client, key, body)
			row.LatencyMS = time.Since(start).Milliseconds()
			cancel()
			var apiErr *APIError
			if errors.As(err, &apiErr) && (apiErr.Status == http.StatusUnauthorized || apiErr.Status == http.StatusUnprocessableEntity) {
				return fmt.Errorf("%s: Jev refused the request: %v", row.File, err)
			}
			if err != nil {
				row.Err = err.Error()
			}
			if cassette != "" && err == nil {
				if err := os.MkdirAll(cassette, 0o755); err != nil {
					return err
				}
				if err := os.WriteFile(filepath.Join(cassette, row.Hash+".json"), data, 0o644); err != nil {
					return err
				}
			}
		}
		m.Entries = append(m.Entries, manifestEntry{File: row.File, Request: row.Hash, LatencyMS: row.LatencyMS, Error: row.Err})
		if row.Err == "" {
			if row.Answer, err = jev.ParseState(data); err != nil {
				row.Err = "invalid answer: " + err.Error()
			}
		}
		done = append(done, row)
	}
	if cassette != "" && !replay {
		data, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(cassette, "manifest.json"), append(data, '\n'), 0o644); err != nil {
			return err
		}
	}
	if replay {
		if old, err := readManifestFile(cassette); err == nil {
			m.Date = old.Date
		}
	}
	reportState(w, m, done)
	if replay {
		var stale []string
		for _, e := range recorded {
			stale = append(stale, e.File)
		}
		slices.Sort(stale)
		drift(w, notRecorded, stale)
	}
	return nil
}

func reportState(w io.Writer, m manifest, rows []StateRow) {
	fmt.Fprintf(w, "Model `%s`, run %s, state prompt `%s`, %d screens.\n\n", m.Model, m.Date, m.Prompt, len(rows))
	fmt.Fprintln(w, "| # | File | Fixture state | Jev | Match | Confidence | Latency ms |")
	fmt.Fprintln(w, "| ---: | --- | --- | --- | --- | ---: | ---: |")
	scored, matched, kept, keptMatch, errs := 0, 0, 0, 0, 0
	var lat []int64
	for i, r := range rows {
		jevState, conf := "error: "+r.Err, "-"
		if r.Err == "" {
			jevState, conf = r.Answer.Choice, fmt.Sprintf("%.3f", r.Answer.Confidence)
		} else {
			errs++
		}
		ok := "unscored"
		if r.Scored() {
			scored++
			ok = "no"
			if r.Match() {
				ok = "yes"
				matched++
			}
			if r.Err == "" && r.Answer.Confidence >= stateThreshold {
				kept++
				if r.Match() {
					keptMatch++
				}
			}
		}
		lat = append(lat, r.LatencyMS)
		fmt.Fprintf(w, "| %d | `%s` | %s | %s | %s | %s | %d |\n", i+1, short(r.File), r.Label, jevState, ok, conf, r.LatencyMS)
	}
	fmt.Fprintln(w, "\n### Unscored screens (the fixture reads unknown)")
	fmt.Fprintln(w)
	n := 0
	for _, r := range rows {
		if !r.Scored() {
			n++
			if r.Err == "" {
				fmt.Fprintf(w, "- `%s`: Jev %s (%.3f)\n", short(r.File), r.Answer.Choice, r.Answer.Confidence)
			} else {
				fmt.Fprintf(w, "- `%s`: error %s\n", short(r.File), r.Err)
			}
		}
	}
	if n == 0 {
		fmt.Fprintln(w, "None.")
	}
	fmt.Fprintln(w, "\n### Wrong at or above 0.85 on the scored screens")
	fmt.Fprintln(w)
	n = 0
	for _, r := range rows {
		if r.Scored() && r.Err == "" && !r.Match() && r.Answer.Confidence >= stateThreshold {
			n++
			fmt.Fprintf(w, "- `%s`: fixture %s, Jev %s (%.3f)\n", short(r.File), r.Label, r.Answer.Choice, r.Answer.Confidence)
		}
	}
	if n == 0 {
		fmt.Fprintln(w, "None.")
	}
	slices.Sort(lat)
	fmt.Fprintln(w, "\n### Summary")
	fmt.Fprintln(w)
	if scored > 0 {
		fmt.Fprintf(w, "- Match on the scored screens: %d/%d (%.1f%%); %d screens unscored.\n", matched, scored, 100*float64(matched)/float64(scored), len(rows)-scored)
	}
	fmt.Fprintf(w, "- Scored answers at or above %.2f: %d, of which %d match.\n", stateThreshold, kept, keptMatch)
	fmt.Fprintf(w, "- Failed requests or invalid answers: %d.\n", errs)
	if len(lat) > 0 {
		fmt.Fprintf(w, "- Latency over %d requests: p50 %d ms, p95 %d ms, min %d ms, max %d ms.\n", len(lat), pct(lat, 50), pct(lat, 95), lat[0], lat[len(lat)-1])
	}
}
