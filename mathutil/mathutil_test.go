package mathutil

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

func TestRound(t *testing.T) {
	cases := []struct {
		in   float64
		want float64
	}{
		{123.4999, 123},
		{123.5, 124},
		{123.999, 124},
		{-123.5, -124},
	}

	for _, c := range cases {
		got := Round(c.in)
		if got != c.want {
			t.Errorf("Round(%f) => %f, want %f", c.in, got, c.want)
		}
	}

}

func TestRoundPlus(t *testing.T) {
	cases := []struct {
		in        float64
		precision int
		want      float64
	}{
		{123.554999, 3, 123.555},
		{123.555555, 3, 123.556},
		{123.558, 2, 123.56},
		{-123.555555, 3, -123.556},
		{123.233333, 2, 123.23},
	}

	for _, c := range cases {
		got := RoundPlus(c.in, c.precision)
		if got != c.want {
			t.Errorf("Round(%f) => %f, want %f", c.in, got, c.want)
		}
	}

}

func TestCeilPlus(t *testing.T) {
	cases := []struct {
		in        float64
		precision int
		want      float64
	}{
		{123.554999, 3, 123.555},
		{123.555555, 3, 123.556},
		{123.558, 2, 123.56},
		{-123.555555, 3, -123.555},
		{123.233333, 2, 123.24},
	}

	for _, c := range cases {
		got := CeilPlus(c.in, c.precision)
		if got != c.want {
			t.Errorf("CeilPlus(%f) => %f, want %f", c.in, got, c.want)
		}
	}
}

func TestIsSignedZero(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"1", false},
		{"0", false},
		{"-1", false},
		{"-0", true},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			f, err := strconv.ParseFloat(tc.in, 64)
			if err != nil {
				t.Fatal(err)
			}

			out := IsSignedZero(f)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestMebi(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{500, "500.0B"},
		{1023, "1023.0B"},
		{1024, "1.0KiB"},
		{1424, "1.4KiB"},
		{152310, "148.7KiB"},
		{1024 * 1190, "1.2MiB"},
		{(math.Pow(1024, 5) * 3) + (math.Pow(1024, 4) * 400), "3.4PiB"},
		{(math.Pow(1024, 6) * 3) + (math.Pow(1024, 5) * 400), "3472.0PiB"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out := Byte(tc.in).String()
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestByte(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{500, "500.0B"},
		{1023, "1.0KB"},
		{1000, "1.0KB"},
		{1424, "1.4KB"},
		{152310, "152.3KB"},
		{1024 * 1190, "1.2MB"},
		{(math.Pow(1024, 5) * 3) + (math.Pow(1024, 4) * 400), "3.8PB"},
		{(math.Pow(1024, 6) * 3) + (math.Pow(1024, 5) * 400), "3909.1PB"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out := Byte(tc.in).StringAsBytes()
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}
