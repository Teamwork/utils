package errorutil

import (
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
)

func TestFilterExclude(t *testing.T) {
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
}

func TestFilterExcludeAll(t *testing.T) {
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
}

func TestFilterInclude(t *testing.T) {
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
}

func TestFilterNil(t *testing.T) {
	var err error
	err = FilterTrace(err, FilterPattern(FilterTraceExlude, "testing"))
	if err != nil {
		t.Errorf("wrong error: %v", err)
	}
}

func makeErr() error {
	err := zxc()
	return errors.WithStack(err)
}

func zxc() error {
	_, err := ioutil.ReadFile("/var/empty/nonexistent")
	return errors.Wrap(err, "could not read")
}
