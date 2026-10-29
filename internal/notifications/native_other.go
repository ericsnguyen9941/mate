//go:build !darwin

package notifications

import "context"

func sendNative(context.Context, Message) error { return ErrUnsupported }
