package contextutil

import (
	"context"
	"time"
)

// detachedContext is a context that is detached from the parent context.
type detachedContext struct {
	ctx context.Context
}

// Deadline returns the time when work done on behalf of this context
// should be canceled. In a detached context, it will always return false.
func (d detachedContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

// Done returns a nil channel always.
func (d detachedContext) Done() <-chan struct{} {
	return nil
}

// If Done is not yet closed, Err returns nil.
// If Done is closed, Err returns a non-nil error explaining why:
// Canceled if the context was canceled
// or DeadlineExceeded if the context's deadline passed.
// In a detached context, it will always return false.
func (d detachedContext) Err() error {
	return nil
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (d detachedContext) Value(key interface{}) interface{} {
	return d.ctx.Value(key)
}

// DetachContext returns a context that is detached from the parent context.
// This is useful when you want to run a function in a goroutine that will
// outlive the parent context.
func DetachContext(ctx context.Context) context.Context {
	return detachedContext{ctx}
}
