package pkg

import (
	"context"
	"time"
)

// DetachContext returns a context that keeps all the values of its parent context
// but detaches from the cancellation and error handling.
func DetachContext(ctx context.Context) context.Context {
	return detachedContext{ctx}
}

type detachedContext struct {
	parent context.Context
}

func (v detachedContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (v detachedContext) Done() <-chan struct{} {
	return nil
}

func (v detachedContext) Err() error {
	return nil
}

func (v detachedContext) Value(key any) any {
	return v.parent.Value(key)
}
