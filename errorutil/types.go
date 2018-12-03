package errorutil

import "github.com/pkg/errors"

type (
	causer interface {
		Cause() error
	}

	stackTracer interface {
		StackTrace() errors.StackTrace
	}

	withStack struct {
		err   error
		stack errors.StackTrace
	}
)

func (w *withStack) Cause() error                  { return w.err }
func (w *withStack) StackTrace() errors.StackTrace { return w.stack }

func (w *withStack) Error() string {
	if w.err == nil {
		return ""
	}
	return w.err.Error()
}
