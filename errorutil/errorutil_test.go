package errorutil

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime"
	"testing"

	"github.com/pkg/errors"
)

func TestFilterTrace(t *testing.T) {
	zxc := func() error {
		_, err := ioutil.ReadFile("/var/empty/nonexistent")
		return errors.Wrap(err, "could not read")
	}
	makeErr := func() error {
		err := zxc()
		return errors.WithStack(err)
	}

	t.Run("exclude", func(t *testing.T) {
		err := makeErr()
		err = errors.Wrap(err, "w00t")
		err = errors.Wrap(err, "context")
		err = FilterTrace(err, FilterPattern(FilterTraceExlude,
			"testing",
			"re:.*github.com/teamwork/utils/.*",
		))

		tErr, _ := err.(stackTracer)
		if len(tErr.StackTrace()) != 1 {
			t.Errorf("wrong length for stack trace: %d; wanted 1", len(tErr.StackTrace()))
			for _, f := range tErr.StackTrace() {
				t.Logf("%+v\n", f)
			}
		}
	})

	t.Run("exclude-all", func(t *testing.T) {
		err := makeErr()
		err = errors.Wrap(err, "w00t")
		err = errors.Wrap(err, "context")
		err = FilterTrace(err, FilterPattern(FilterTraceExlude, "re:.*"))

		tErr, _ := err.(stackTracer)
		if len(tErr.StackTrace()) != 3 {
			t.Errorf("wrong length for stack trace: %d; wanted 3", len(tErr.StackTrace()))
			for _, f := range tErr.StackTrace() {
				t.Logf("%+v\n", f)
			}
		}
	})

	t.Run("include", func(t *testing.T) {
		err := makeErr()
		err = errors.Wrap(err, "w00t")
		err = errors.Wrap(err, "context")
		err = FilterTrace(err, FilterPattern(FilterTraceInclude,
			"re:.*github.com/teamwork/utils/.*"))

		tErr, _ := err.(stackTracer)
		if len(tErr.StackTrace()) != 1 {
			t.Errorf("wrong length for stack trace: %d; wanted 1", len(tErr.StackTrace()))
			for _, f := range tErr.StackTrace() {
				t.Logf("%+v\n", f)
			}
		}
	})

	t.Run("nil", func(t *testing.T) {
		var err error
		err = FilterTrace(err, FilterPattern(FilterTraceExlude, "testing"))
		if err != nil {
			t.Errorf("wrong error: %v", err)
		}
	})
}

func TestEarliestStackTracer(t *testing.T) {
	stack := errors.WithStack(fmt.Errorf("w00t"))
	deeper := errors.WithStack(stack)

	tests := []struct {
		in   error
		want error
	}{
		{nil, nil},
		{fmt.Errorf("w00t"), nil},
		{stack, stack},
		{deeper, stack},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := EarliestStackTracer(tt.in)
			if out != tt.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}

func TestAddStackTrace(t *testing.T) {
	err := fmt.Errorf("w00t")
	stack := errors.WithStack(err)
	errStack := &withStack{
		err:   err,
		stack: stack.(stackTracer).StackTrace()[1:],
	}
	emptyStack := &withStack{err: err, stack: []errors.Frame{}}

	tests := []struct {
		in       error
		inIgnore string
		want     error
	}{
		{nil, "", nil},
		{stack, "", stack},
		{err, "", errStack},
		{err, runtime.GOROOT(), emptyStack},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := AddStackTrace(tt.in, tt.inIgnore)
			if !reflect.DeepEqual(tt.want, out) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tt.want)
			}
		})
	}
}
