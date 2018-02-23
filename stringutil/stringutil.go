// Package stringutil adds functions for working with strings.
package stringutil

import "regexp"

// Left returns the "n" left characters of the string.
//
// If the string is shorter than "n" it will return the first "n" characters of
// the string with "…" appended. Otherwise the entire string is returned as-is.
func Left(s string, n int) string {
	if n < 0 {
		n = 0
	}
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

// RemoveUnprintable removes unprintable characters (0 to 31 ASCII) from a string.
func RemoveUnprintable(s string) (string, error) {
	r, err := regexp.Compile("[\x00-\x1F]")
	if err != nil {
		return s, err
	}
	return r.ReplaceAllString(s, ""), nil
}
