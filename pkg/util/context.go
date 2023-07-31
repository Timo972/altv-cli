package util

import (
	"context"
	"time"
)

// ContextWithOptionalTimeout returns a context with a timeout if the timeout (in s) is greater than 0, otherwise it returns a cancelable context.
func ContextWithOptionalTimeout(ctx context.Context, timeout int) (context.Context, context.CancelFunc) {
	if timeout > 0 {
		return context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	}
	return context.WithCancel(ctx)
}
