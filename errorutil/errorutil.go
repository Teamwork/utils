// Package errorutil provides functions to work with errors.
package errorutil // import "github.com/teamwork/utils/v2/errorutil"

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// Modes for FilterPatterns.
const (
	FilterTraceExclude = 0 // Exclude the paths that match.
	FilterTraceInclude = 1 // Include only the paths that match.
)

// Patterns for filtering error traces.
type Patterns struct {
	ret     bool
	files   []string
	pkgs    []string
	matches []string
	regexps []*regexp.Regexp
}

// FilterPattern compiles filter patterns for FilterTrace().
// Note: After go 1.13 this can include go run time data, excluding "re:.*runtime/*.s" is recommended.
//
// Frames are filtered according to the mode; with FilterTraceExclude all frames
// are included except those that match the given patterns. With
// FilterTraceInclude all frames are excluded except those that match one of the
// patterns.
//
// Paths starting with re: are treated as a regular expression.
//
// Paths starting with match: are matched with filepath.Match()
//
// Paths ending with .go are matches against the full file path (i.e.
// /home/martin/go/src/.../file.go).
//
// Anything else is matches against the package path (i.e. github.com/foo/bar).
func FilterPattern(mode int, paths ...string) *Patterns {
	var pat Patterns
	switch mode {
	case FilterTraceExclude:
		pat.ret = true
	case FilterTraceInclude:
		pat.ret = false
	default:
		panic(fmt.Sprintf("FilterPattern: invalid mode: %q", mode))
	}

	for _, p := range paths {
		switch {
		case strings.HasPrefix(p, "match:"):
			// Make sure pattern isn't malformed.
			_, err := filepath.Match(p, "")
			if err != nil {
				panic(fmt.Sprintf("FilterPattern: invalid match pattern: %s", err))
			}

			pat.matches = append(pat.matches, p[6:])
		case strings.HasPrefix(p, "re:"):
			pat.regexps = append(pat.regexps, regexp.MustCompile(p[3:]))
		case strings.HasSuffix(p, ".go"):
			pat.files = append(pat.files, p)
		default:
			pat.pkgs = append(pat.pkgs, p)
		}
	}

	return &pat
}

// Match a file path.
func (p Patterns) Match(pc uintptr) bool {
	fn := runtime.FuncForPC(pc)
	file, _ := fn.FileLine(pc)

	for _, f := range p.files {
		if file == f {
			return p.ret
		}
	}

	if len(p.pkgs) > 0 {
		// Get package name.
		pkg := fn.Name()
		s := strings.LastIndex(pkg, "/")
		if s < 0 {
			s = 0
		}
		if d := strings.Index(pkg[s:], "."); d > -1 {
			pkg = pkg[:s+d]
		}
		if v := strings.Index(pkg, "/vendor/"); v > -1 {
			pkg = pkg[v+8:]
		}
		for _, d := range p.pkgs {
			if strings.HasPrefix(pkg, d) {
				return p.ret
			}
		}
	}

	for _, m := range p.matches {
		if ok, err := filepath.Match(m, file); ok && err == nil {
			return p.ret
		}
	}

	for _, r := range p.regexps {
		if r.MatchString(file) {
			return p.ret
		}
	}

	return !p.ret
}

// FilterTrace removes unneeded stack traces from an error.
func FilterTrace(err error, p *Patterns) error {
	tErr, ok := err.(stackTracer)
	if !ok {
		return err
	}

	var frames errors.StackTrace
	for _, frame := range tErr.StackTrace() {
		if !p.Match(uintptr(frame) - 1) {
			frames = append(frames, frame)
		}
	}

	// Keep original stack if we filtered everything, because that's not likely
	// going to be useful.
	if len(frames) == 0 {
		_, _ = fmt.Fprintf(os.Stderr,
			"WARNING: errorutil.FilterTrace: all stack frames filtered; keeping full trace\n")
		frames = tErr.StackTrace()
	}

	return &withStack{
		err:   err,
		stack: frames,
	}
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type withStack struct {
	err   error
	stack errors.StackTrace
}

func (w *withStack) Cause() error                  { return w.err }
func (w *withStack) StackTrace() errors.StackTrace { return w.stack }

func (w *withStack) Error() string {
	if w.err == nil {
		return ""
	}
	return w.err.Error()
}
