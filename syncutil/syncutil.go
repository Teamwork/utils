// Package syncutil adds functions for synchronization.
package syncutil // import "github.com/teamwork/utils/syncutil"

import (
	"context"
	"sync"
)

// Wait for a sync.WaitGroup with support for timeout/cancellations from
// context.
func Wait(ctx context.Context, wg *sync.WaitGroup) error {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		wg.Wait()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return nil
	}
}
