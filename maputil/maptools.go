// Package maputil provides a set if functions for working with maps.
package maputil // import "github.com/teamwork/utils/maputil"

// Reverse the keys and values of a map.
func Reverse(m map[string]string) map[string]string {
	n := make(map[string]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}
