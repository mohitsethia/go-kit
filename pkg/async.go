package pkg

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
)

// Async ...
type Async interface {
	RunAsync(
		ctx context.Context,
		name string,
		callback func(context.Context) error,
	)
	Wait()
}

type async struct {
	wg sync.WaitGroup
}

// NewAsync ...
func NewAsync() Async {
	return &async{}
}

// RunAsync ...
func (ae *async) RunAsync(
	ctx context.Context,
	name string,
	cb func(context.Context) error,
) {
	ae.wg.Add(1)
	go func(ctx context.Context) {
		defer ae.wg.Done()
		defer CapturePanic(func(err error) {
			if err == nil {
				return
			}
			log.Errorf("[ASYNC][%s] panic: %v", name, err)
		})
		if err := cb(ctx); err != nil {
			log.Errorf("[ASYNC][%s] error: %v", name, err)
		}
	}(ctx)
}

// Wait ...
func (ae *async) Wait() {
	ae.wg.Wait()
}

// CapturePanic should be used with defer to be able to recover
// from panic gracefully
// NOTE: this function can't be refactored in something like
// `func CapturePanic() error` and used from any deferred function to recover panic,
// because in this case it will loose panic context and won't be able to recover
func CapturePanic(cb func(error)) {
	if r := recover(); r == nil {
		cb(nil)
	} else {
		cb(fmt.Errorf("panic: %+v, stack: %s",
			r, bytes.ReplaceAll(debug.Stack(), []byte("\n\t"), []byte(" -> "))))
	}
}
