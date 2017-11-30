// Package stringutil contains functions for working with strings.
package stringutil

import "strings"

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

// Indent every line in a string with 'n' repetitions of 'indent'.
//
// Use a negative number to remove indentation. If n is 0 it will de-indent to
// match the first line.
func Indent(str, indent string, n int) string {
	if n == 0 {
		for i := 0; i < len(str); i++ {
			switch str[i] {
			case '\n':
				// Do nothing
			case '\t': // TODO: support space
				n++
			default:
				break
			}
		}
	}

	endNl := strings.HasSuffix(str, "\n")

	r := ""
	if n < 0 {
		for _, line := range strings.Split(str, "\n") {
			r += strings.Replace(line, indent, "", -n) + "\n"
		}
	} else {
		for _, line := range strings.Split(str, "\n") {
			r += strings.Repeat(indent, n) + line + "\n"
		}
	}

	if endNl {
		return r
	}
	return r[:len(r)-1]
}
