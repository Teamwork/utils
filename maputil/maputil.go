// Package maputil provides a set if functions for working with maps.
package maputil // import "github.com/teamwork/utils/v2/maputil"

// Swap the keys and values of a map.
func Swap[T comparable, V comparable](m map[T]V) map[V]T {
	n := make(map[V]T)
	for k, v := range m {
		n[v] = k
	}

	return n
}
