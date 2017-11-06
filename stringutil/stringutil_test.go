package stringutil

import (
	"fmt"
	"testing"
)

func TestLeft(t *testing.T) {
	cases := []struct {
		in   string
		n    int
		want string
	}{
		{"Hello", 100, "Hello"},
		{"Hello", 1, "H…"},
		{"Hello", 5, "Hello"},
		{"Hello", 4, "Hell…"},
		{"Hello", 0, "…"},
		{"Hello", -2, "…"},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := Left(tc.in, tc.n)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}
