package notifications

import (
	"context"
	"errors"
)

// Native delivers through macOS Notification Center under Mate's own identity.
// The companion is compiled once on macOS; other platforms return an error.
type Native struct{}

func (Native) Send(ctx context.Context, m Message) error { return sendNative(ctx, m) }

var ErrUnsupported = errors.New("native Mate notifications require macOS")
