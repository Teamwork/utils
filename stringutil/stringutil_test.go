package stringutil

import (
	"fmt"
	"strings"
	"testing"
)

func TestTruncate(t *testing.T) {
	cases := []struct {
		in   string
		n    int
		want string
	}{
		{"Hello", 100, "Hello"},
		{"Hello", 1, "H"},
		{"Hello", 5, "Hello"},
		{"Hello", 4, "Hell"},
		{"Hello", 0, ""},
		{"Hello", -2, ""},
		{"汉语漢語", 1, "汉"},
		{"汉语漢語", 3, "汉语漢"},
		{"汉语漢語", 4, "汉语漢語"},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := Truncate(tc.in, tc.n)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

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
		{"汉语漢語", 1, "汉…"},
		{"汉语漢語", 3, "汉语漢…"},
		{"汉语漢語", 4, "汉语漢語"},
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

func TestUpperFirst(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"hello", "Hello"},
		{"helloWorld", "HelloWorld"},
		{"h", "H"},
		{"hh", "Hh"},
		{"ëllo", "Ëllo"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out := UpperFirst(tc.in)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestLowerFirst(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Hello", "hello"},
		{"HelloWorld", "helloWorld"},
		{"H", "h"},
		{"HH", "hH"},
		{"Ëllo", "ëllo"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out := LowerFirst(tc.in)
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

func TestStripWhitespaces(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{{
		name: "it should strip all whitespaces",
		in:   " A test phrase ",
		want: "Atestphrase",
	}, {
		name: "it should return the same string",
		in:   "no_whitespace",
		want: "no_whitespace",
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := StripWhitespaces(tc.in)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func BenchmarkLeft(b *testing.B) {
	text := strings.Repeat("Hello, world, it's a sentences!\n", 200)
	for n := 0; n < b.N; n++ {
		Left(text, 250)
	}
}

func BenchmarkRemoveUnprintable(b *testing.B) {
	text := strings.Repeat("Hello, world, it's a sentences!\n", 20000)
	for n := 0; n < b.N; n++ {
		GetLine(text, 200)
	}
}
