// Package contextutil adds functions for context.
package contextutil // import "github.com/teamwork/utils/v2/contextutil"

import (
	"context"
	"time"
)

// detachedContext is a context that is detached from the parent context.
type detachedContext struct {
	ctx context.Context
}

// Deadline returns a zeroed time and false.
// Refer to `context.Context.Deadline` for an explanation of what `Deadline`
// usually does.
func (d detachedContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

// Done always returns a nil channel.
// Refer to `context.Context.Done` for an explanation of what `Done`
// usually does.
func (d detachedContext) Done() <-chan struct{} {
	return nil
}

// Err always returns nil.
// Refer to `context.Context.Err` for an explanation of what `Err`
// usually does.
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
