// Package dbg contains helper functions useful when debugging programs.
package dbg // import "github.com/teamwork/utils/v2/dbg"

import (
	"fmt"
	"runtime"
)

// Loc gets a location in the stack trace.
//
// Use 0 for the current location; 1 for one up, etc.
func Loc(n int) string {
	_, file, line, ok := runtime.Caller(n + 1)
	if !ok {
		file = "???"
		line = 0
	}

	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	return fmt.Sprintf("%v:%v", file, line)
}
