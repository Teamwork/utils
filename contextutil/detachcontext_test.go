package contextutil_test

import (
	"context"
	"testing"
	"time"

	"github.com/teamwork/utils/v2/contextutil"
)

func TestDetach(t *testing.T) {
	t.Parallel()
	type myCoolKey struct{}

	pCtx := context.WithValue(context.Background(), myCoolKey{}, 56)

	pCtx, cancel := context.WithDeadline(pCtx, time.Now().AddDate(0, 0, -1))
	defer cancel()

	ctx := contextutil.DetachContext(pCtx)

	if pCtx.Err() == nil {
		t.Fatal("expected err")
	}
	if ctx.Err() != nil {
		t.Fatal("expected nil err")
	}
	if pCtx.Done() == nil {
		t.Fatal("expected done channel")
	}
	if ctx.Done() != nil {
		t.Fatal("expected nil done channel")
	}

	if ctx.Value(myCoolKey{}).(int) != 56 {
		t.Fatal("expected to retrieve value 56")
	}
}
