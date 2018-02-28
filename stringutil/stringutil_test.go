package stringutil

import (
	"fmt"
	"strings"
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

func TestRemoveUnprintable(t *testing.T) {
	cases := []struct {
		in      string
		lenLost int
		want    string
	}{
		{"Hello, 世界", 0, "Hello, 世界"},
		{"m", 1, "m"},
		{"m", 0, "m"},
		{" ", 3, " "},
		{"a‎b‏c", 6, "abc"}, // only 2 removed but count as 3 each
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := RemoveUnprintable(tc.in)
			charsRemoved := len(tc.in) - len(out)
			if tc.lenLost != charsRemoved {
				t.Errorf("\ncharsRemoved:  %#v\nwant: %#v\n", charsRemoved, tc.lenLost)
			}
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestGetLine(t *testing.T) {
	cases := []struct {
		in   string
		line int
		want string
	}{
		{"Hello", 1, "Hello"},
		{"Hello", 2, ""},
		{"Hello\nworld", 1, "Hello"},
		{"Hello\nworld", 2, "world"},
		{"Hello\nworld", 3, ""},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out := GetLine(tc.in, tc.line)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func BenchmarkRemoveUnprintable(b *testing.B) {
	text := strings.Repeat("Hello, world, it's a sentences!\n", 20000)
	for n := 0; n < b.N; n++ {
		GetLine(text, 200)
	}
}
