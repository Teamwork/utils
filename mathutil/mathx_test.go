package mathutil

import (
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
	}

	for _, c := range cases {
		got := RoundPlus(c.in, c.precision)
		if got != c.want {
			t.Errorf("Round(%f) => %f, want %f", c.in, got, c.want)
		}
	}

}
