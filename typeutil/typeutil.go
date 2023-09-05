// Package typeutil adds functions for types.
package typeutil // import "github.com/teamwork/utils/v2/typeutil"

// Default returns `val` if it is not zero, otherwise returns
// `def`.
//
// v := Default("", "hello")      // return "hello"
// v := Default("world", "hello") // return "world"
func Default[T comparable](val, def T) T {
	if val == *new(T) {
		return def
	}

	return val
}
