package syncutil

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	t.Run("cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := Wait(ctx, &wg)
		if err != context.Canceled {
			t.Errorf("wrong error: %v", err)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := Wait(ctx, &wg)
		if err != context.DeadlineExceeded {
			t.Errorf("wrong error: %v", err)
		}
	})

	t.Run("finish", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		wg.Done()
		wg.Done()

		err := Wait(ctx, &wg)
		if err != nil {
			t.Errorf("wrong error: %v", err)
		}
	})
}
