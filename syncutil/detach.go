package syncutil

import (
	"context"
	"time"
)

// detachedContext is a context that is detached from the parent context.
type detachedContext struct {
	ctx context.Context
}

// Deadline returns the time when work done on behalf of this context
func (d detachedContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

// Done returns a nil channel always.
func (d detachedContext) Done() <-chan struct{} {
	return nil
}

// Err returns nil always.
func (d detachedContext) Err() error {
	return nil
}

// Value returns the value associated with this context for key, or nil if no
func (d detachedContext) Value(key interface{}) interface{} {
	return d.ctx.Value(key)
}

// DetachContext returns a context that is detached from the parent context.
// This is useful when you want to run a function in a goroutine that will
// outlive the parent context.
func DetachContext(ctx context.Context) context.Context {
	return detachedContext{ctx}
}
