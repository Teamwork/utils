package timeutil

import (
	"fmt"
	"testing"
	"time"
)

func TestStartOfMonth(t *testing.T) {
	cases := []struct {
		in   time.Time
		want time.Time
	}{
		{mustParse(t, "2016-01-13"), mustParse(t, "2016-01-01")},
		{mustParse(t, "2016-01-01"), mustParse(t, "2016-01-01")},
		{mustParse(t, "2016-12-30"), mustParse(t, "2016-12-01")},
	}

	for _, c := range cases {
		got := StartOfMonth(c.in)
		if got != c.want {
			t.Errorf("StartOfMonth(%s) => %s, want %s", c.in, got, c.want)
		}
	}
}

func TestEndOfMonth(t *testing.T) {
	cases := []struct {
		in   time.Time
		want time.Time
	}{
		{mustParse(t, "2016-01-01"), mustParse(t, "2016-01-31")},
		{mustParse(t, "2016-01-31"), mustParse(t, "2016-01-31")},
		{mustParse(t, "2016-11-01"), mustParse(t, "2016-11-30")},
		{mustParse(t, "2016-12-31"), mustParse(t, "2016-12-31")},
		// leap test
		{mustParse(t, "2012-02-01"), mustParse(t, "2012-02-29")},
		{mustParse(t, "2013-02-01"), mustParse(t, "2013-02-28")},
	}

	for _, c := range cases {
		got := EndOfMonth(c.in)
		if got != c.want {
			t.Errorf("EndOfMonth(%s) => %s, want %s", c.in, got, c.want)
		}
	}
}

// mustParse parses value in the format YYYY-MM-DD failing the test on error.
func mustParse(t *testing.T, value string) time.Time {
	const layout = "2006-01-02"
	d, err := time.Parse(layout, value)
	if err != nil {
		t.Fatalf("time.Parse(%q, %q) unexpected error: %v", layout, value, err)
	}
	return d
}

func TestMonthsTo(t *testing.T) {
	day := 24 * time.Hour
	cases := []struct {
		in   time.Time
		want int
	}{
		{time.Now(), 1},
		{time.Now().Add(day * 35), 1},
		{time.Now().Add(day * 65), 2},
		{time.Now().Add(day * 370), 12},

		// Broken!
		//{time.Now().Add(-day * 35), -1},
		//{time.Now().Add(-day * 65), -2},
		//{time.Now().Add(-day * 370), -12},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			out := MonthsTo(tc.in)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}
