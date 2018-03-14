package timeutil

import (
	"fmt"
	"testing"
	"time"
)

func TestStartOfDay(t *testing.T) {
	cases := []struct {
		in   time.Time
		want time.Time
	}{
		{mustParse(t, "2016-01-01T00:00:00Z"), mustParse(t, "2016-01-01T00:00:00Z")},
		{mustParse(t, "2016-01-13T00:00:00Z"), mustParse(t, "2016-01-13T00:00:00Z")},
		{mustParse(t, "2016-01-01T12:34:56Z"), mustParse(t, "2016-01-01T00:00:00Z")},
		{mustParse(t, "2016-12-30T23:59:59Z"), mustParse(t, "2016-12-30T00:00:00Z")},
	}

	for _, c := range cases {
		got := StartOfDay(c.in)
		if got != c.want {
			t.Errorf("StartOfDay(%s) =>\n%s, want %s", c.in, got, c.want)
		}
	}
}

func TestEndOfDay(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Dublin")
	if err != nil {
		t.Fatalf("LoadLocation: Europe/Dublin: %v", err)
	}

	cases := []struct {
		in   time.Time
		want time.Time
	}{
		{mustParse(t, "2016-01-01T00:00:00Z"), mustParse(t, "2016-01-01T23:59:59.999999999Z")},
		{mustParse(t, "2016-01-13T00:00:00Z"), mustParse(t, "2016-01-13T23:59:59.999999999Z")},
		{mustParse(t, "2016-01-01T12:34:56Z"), mustParse(t, "2016-01-01T23:59:59.999999999Z")},
		{mustParse(t, "2016-12-30T23:59:59Z"), mustParse(t, "2016-12-30T23:59:59.999999999Z")},
		{ // dst
			mustParse(t, "2018-03-25T00:00:00Z").In(loc),
			mustParse(t, "2018-03-25T23:59:59.999999999+01:00"),
		},
	}

	for _, c := range cases {
		got := EndOfDay(c.in)
		if got.UTC() != c.want.UTC() {
			t.Errorf("EndOfDay(%s) =>\n%s, want %s", c.in, got, c.want)
		}
	}
}

func TestTomorrow(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Dublin")
	if err != nil {
		t.Fatalf("LoadLocation: Europe/Dublin: %v", err)
	}

	cases := []struct {
		in   time.Time
		want time.Time
	}{
		{mustParse(t, "2016-01-01T00:00:00Z"), mustParse(t, "2016-01-02T00:00:00Z")},
		{mustParse(t, "2016-01-13T00:00:00Z"), mustParse(t, "2016-01-14T00:00:00Z")},
		{mustParse(t, "2016-01-01T12:34:56Z"), mustParse(t, "2016-01-02T00:00:00Z")},
		{mustParse(t, "2016-12-30T23:59:59Z"), mustParse(t, "2016-12-31T00:00:00Z")},
		{mustParse(t, "2016-12-31T00:00:00Z"), mustParse(t, "2017-01-01T00:00:00Z")},
		{ // dst
			mustParse(t, "2018-03-25T00:00:00Z").In(loc),
			mustParse(t, "2018-03-26T00:00:00+01:00"),
		},
	}

	for _, c := range cases {
		got := Tomorrow(c.in)
		if got.UTC() != c.want.UTC() {
			t.Errorf("Tomorrow(%s) =>\n%s, want %s", c.in, got, c.want)
		}
	}
}

func TestStartOfWeek(t *testing.T) {
	cases := []struct {
		in   time.Time
		want time.Time
	}{
		{mustParse(t, "2018-01-01T00:00:00Z"), mustParse(t, "2018-01-01T00:00:00Z")},
		{mustParse(t, "2018-01-02T00:00:00Z"), mustParse(t, "2018-01-01T00:00:00Z")},
		{mustParse(t, "2018-01-03T00:00:00Z"), mustParse(t, "2018-01-01T00:00:00Z")},
		{mustParse(t, "2018-01-04T00:00:00Z"), mustParse(t, "2018-01-01T00:00:00Z")},
		{mustParse(t, "2018-01-05T00:00:00Z"), mustParse(t, "2018-01-01T00:00:00Z")},
		{mustParse(t, "2018-01-06T00:00:00Z"), mustParse(t, "2018-01-01T00:00:00Z")},
		{mustParse(t, "2018-01-07T00:00:00Z"), mustParse(t, "2018-01-01T00:00:00Z")},

		{mustParse(t, "2018-01-08T12:34:56Z"), mustParse(t, "2018-01-08T00:00:00Z")},
		{mustParse(t, "2018-01-09T00:00:00Z"), mustParse(t, "2018-01-08T00:00:00Z")},
		{mustParse(t, "2018-01-10T00:00:00Z"), mustParse(t, "2018-01-08T00:00:00Z")},
		{mustParse(t, "2018-01-11T00:00:00Z"), mustParse(t, "2018-01-08T00:00:00Z")},
		{mustParse(t, "2018-01-12T00:00:00Z"), mustParse(t, "2018-01-08T00:00:00Z")},
		{mustParse(t, "2018-01-13T00:00:00Z"), mustParse(t, "2018-01-08T00:00:00Z")},
		{mustParse(t, "2018-01-14T00:00:00Z"), mustParse(t, "2018-01-08T00:00:00Z")},
	}

	for _, c := range cases {
		got := StartOfWeek(c.in)
		if got != c.want {
			t.Errorf("StartOfWeek(%s) =>\n%s, want %s", c.in, got, c.want)
		}
	}
}

func TestEndOfWeek(t *testing.T) {
	cases := []struct {
		in   time.Time
		want time.Time
	}{
		{mustParse(t, "2018-01-01T00:00:00Z"), mustParse(t, "2018-01-07T23:59:59.999999999Z")},
		{mustParse(t, "2018-01-02T00:00:00Z"), mustParse(t, "2018-01-07T23:59:59.999999999Z")},
		{mustParse(t, "2018-01-03T00:00:00Z"), mustParse(t, "2018-01-07T23:59:59.999999999Z")},
		{mustParse(t, "2018-01-04T00:00:00Z"), mustParse(t, "2018-01-07T23:59:59.999999999Z")},
		{mustParse(t, "2018-01-05T00:00:00Z"), mustParse(t, "2018-01-07T23:59:59.999999999Z")},
		{mustParse(t, "2018-01-06T00:00:00Z"), mustParse(t, "2018-01-07T23:59:59.999999999Z")},
		{mustParse(t, "2018-01-07T00:00:00Z"), mustParse(t, "2018-01-07T23:59:59.999999999Z")},

		{mustParse(t, "2018-01-08T12:34:56Z"), mustParse(t, "2018-01-14T23:59:59.999999999Z")},
		{mustParse(t, "2018-01-09T00:00:00Z"), mustParse(t, "2018-01-14T23:59:59.999999999Z")},
		{mustParse(t, "2018-01-10T00:00:00Z"), mustParse(t, "2018-01-14T23:59:59.999999999Z")},
		{mustParse(t, "2018-01-11T00:00:00Z"), mustParse(t, "2018-01-14T23:59:59.999999999Z")},
		{mustParse(t, "2018-01-12T00:00:00Z"), mustParse(t, "2018-01-14T23:59:59.999999999Z")},
		{mustParse(t, "2018-01-13T00:00:00Z"), mustParse(t, "2018-01-14T23:59:59.999999999Z")},
		{mustParse(t, "2018-01-14T00:00:00Z"), mustParse(t, "2018-01-14T23:59:59.999999999Z")},
	}

	for _, c := range cases {
		got := EndOfWeek(c.in)
		if got != c.want {
			t.Errorf("EndOfWeek(%s) =>\n%s, want %s", c.in, got, c.want)
		}
	}
}

func TestStartOfMonth(t *testing.T) {
	cases := []struct {
		in   time.Time
		want time.Time
	}{
		{mustParse(t, "2016-01-13T00:00:00Z"), mustParse(t, "2016-01-01T00:00:00Z")},
		{mustParse(t, "2016-01-01T12:34:56Z"), mustParse(t, "2016-01-01T00:00:00Z")},
		{mustParse(t, "2016-12-30T00:00:00Z"), mustParse(t, "2016-12-01T00:00:00Z")},
	}

	for _, c := range cases {
		got := StartOfMonth(c.in)
		if got != c.want {
			t.Errorf("StartOfMonth(%s) =>\n%s, want %s", c.in, got, c.want)
		}
	}
}

func TestEndOfMonth(t *testing.T) {
	cases := []struct {
		in   time.Time
		want time.Time
	}{
		{mustParse(t, "2016-01-01T00:00:00Z"), mustParse(t, "2016-01-31T23:59:59.999999999Z")},
		{mustParse(t, "2016-01-31T00:00:00Z"), mustParse(t, "2016-01-31T23:59:59.999999999Z")},
		{mustParse(t, "2016-11-01T00:00:00Z"), mustParse(t, "2016-11-30T23:59:59.999999999Z")},
		{mustParse(t, "2016-12-31T00:00:00Z"), mustParse(t, "2016-12-31T23:59:59.999999999Z")},
		// leap test
		{mustParse(t, "2012-02-01T00:00:00Z"), mustParse(t, "2012-02-29T23:59:59.999999999Z")},
		{mustParse(t, "2013-02-01T00:00:00Z"), mustParse(t, "2013-02-28T23:59:59.999999999Z")},
	}

	for _, c := range cases {
		got := EndOfMonth(c.in)
		if got != c.want {
			t.Errorf("EndOfMonth(%s) =>\n%s, want %s", c.in, got, c.want)
		}
	}
}

func TestStartOfQuarter(t *testing.T) {
	cases := []struct {
		in   time.Time
		want time.Time
	}{
		// q1
		{mustParse(t, "2016-01-01T00:00:00Z"), mustParse(t, "2016-01-01T00:00:00Z")},
		{mustParse(t, "2016-03-31T23:00:00Z"), mustParse(t, "2016-01-01T00:00:00Z")},

		// q2
		{mustParse(t, "2016-04-01T00:00:00Z"), mustParse(t, "2016-04-01T00:00:00Z")},
		{mustParse(t, "2016-06-30T00:00:00Z"), mustParse(t, "2016-04-01T00:00:00Z")},

		// q3
		{mustParse(t, "2016-07-01T00:00:00Z"), mustParse(t, "2016-07-01T00:00:00Z")},
		{mustParse(t, "2016-09-30T00:00:00Z"), mustParse(t, "2016-07-01T00:00:00Z")},

		// q4
		{mustParse(t, "2016-10-01T00:00:00Z"), mustParse(t, "2016-10-01T00:00:00Z")},
		{mustParse(t, "2016-12-31T00:00:00Z"), mustParse(t, "2016-10-01T00:00:00Z")},

		// leap test
		{mustParse(t, "2012-02-01T00:00:00Z"), mustParse(t, "2012-01-01T00:00:00Z")},
		{mustParse(t, "2013-02-01T00:00:00Z"), mustParse(t, "2013-01-01T00:00:00Z")},
	}

	for _, c := range cases {
		got := StartOfQuarter(c.in)
		if got != c.want {
			t.Errorf("StartOfQuarter(%s) =>\n%s, want %s", c.in, got, c.want)
		}
	}
}

func TestEndOfQuarter(t *testing.T) {
	cases := []struct {
		in   time.Time
		want time.Time
	}{
		// q1
		{mustParse(t, "2016-01-01T00:00:00Z"), mustParse(t, "2016-03-31T23:59:59.999999999Z")},
		{mustParse(t, "2016-03-31T00:00:00Z"), mustParse(t, "2016-03-31T23:59:59.999999999Z")},

		// q2
		{mustParse(t, "2016-04-01T00:00:00Z"), mustParse(t, "2016-06-30T23:59:59.999999999Z")},
		{mustParse(t, "2016-06-30T00:00:00Z"), mustParse(t, "2016-06-30T23:59:59.999999999Z")},

		// q3
		{mustParse(t, "2016-07-01T00:00:00Z"), mustParse(t, "2016-09-30T23:59:59.999999999Z")},
		{mustParse(t, "2016-09-30T00:00:00Z"), mustParse(t, "2016-09-30T23:59:59.999999999Z")},

		// q4
		{mustParse(t, "2016-10-01T00:00:00Z"), mustParse(t, "2016-12-31T23:59:59.999999999Z")},
		{mustParse(t, "2016-12-31T00:00:00Z"), mustParse(t, "2016-12-31T23:59:59.999999999Z")},

		// leap test
		{mustParse(t, "2012-02-01T00:00:00Z"), mustParse(t, "2012-03-31T23:59:59.999999999Z")},
		{mustParse(t, "2013-02-01T00:00:00Z"), mustParse(t, "2013-03-31T23:59:59.999999999Z")},
	}

	for _, c := range cases {
		got := EndOfQuarter(c.in)
		if got != c.want {
			t.Errorf("EndOfQuarter(%s) =>\n%s, want %s", c.in, got, c.want)
		}
	}
}

// mustParse parses value in the format time.RFC3339Nano failing the test on error.
func mustParse(t *testing.T, value string) time.Time {
	const layout = time.RFC3339Nano
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
