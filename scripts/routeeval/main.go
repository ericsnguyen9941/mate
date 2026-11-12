// Command routeeval replays crews' hand-backs through the hand-back router
// (internal/route, CON-274) and prints one JSON decision per crew, so a
// project's history can be scored before `autoroute` is turned on.
//
//	go run ./scripts/routeeval -workspace <dir> -project <name> <crew>...
//
// It asks Jev for real, one request per crew that reaches it, with the key
// `.mate/.env` names. It acts on nothing. A crew's meta is read as it is
// now, so what was only true after its hand-back (its pull request merged,
// a review started) is cleared before deciding: the replay asks what the
// router would have decided at the hand-back.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nguyenngocanh94/mate/internal/notice"
	"github.com/nguyenngocanh94/mate/internal/route"
	"github.com/nguyenngocanh94/mate/internal/store"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "routeeval:", err)
		os.Exit(1)
	}
}

func run() error {
	wsFlag := flag.String("workspace", "", "workspace directory")
	projectFlag := flag.String("project", "", "project name")
	review := flag.String("review", store.ReviewOff, "the review setting to decide under")
	flag.Parse()
	ws, err := store.Open(*wsFlag)
	if err != nil {
		return err
	}
	cfg, err := ws.LoadProject(*projectFlag)
	if err != nil {
		return err
	}
	cfg.Review = *review
	key, err := jevKey(ws)
	if err != nil {
		return err
	}
	r := route.New(ws, nil, route.Exec{}, nil)
	judge := notice.New(key)
	enc := json.NewEncoder(os.Stdout)
	for _, crew := range flag.Args() {
		in := r.Inspect(*projectFlag, crew, cfg)
		in.PRMerged, in.Reviewed = false, false
		if err := enc.Encode(route.Decide(context.Background(), in, judge)); err != nil {
			return err
		}
	}
	return nil
}

func jevKey(ws *store.Workspace) (string, error) {
	env, err := ws.LoadEnv()
	if err != nil {
		return "", err
	}
	path := env["MATE_JEV_API_KEY_FILE"]
	if path == "" {
		return "", errors.New("MATE_JEV_API_KEY_FILE is not set in .mate/.env")
	}
	if rest, ok := strings.CutPrefix(path, "~/"); ok {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, rest)
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := io.ReadAll(io.LimitReader(f, 4097))
	return strings.TrimSpace(string(b)), err
}
