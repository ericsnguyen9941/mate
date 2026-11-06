package main

import (
	"context"
	"fmt"
	"io"

	"github.com/nguyenngocanh94/mate/internal/store"
	"github.com/nguyenngocanh94/mate/internal/usage/why"
)

// writeLearningUsageLine prints why.RecordAtClose's line after a close
// report on the command line.
func writeLearningUsageLine(stdout io.Writer, ctx context.Context, w *store.Workspace, project, crew, event string) {
	fmt.Fprintln(stdout, why.RecordAtClose(ctx, w, project, crew, event))
}
