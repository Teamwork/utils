package timeutil

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/teamwork/test/diff"
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

func TestWeeksOnPeriod(t *testing.T) {
	cases := []struct {
		name     string
		period   Period
		optFuncs []PeriodOptionsFunc
		want     []string
	}{
		{
			name: "it should detect 7 weeks between 2 months",
			period: Period{
				Start: time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 12, 11, 0, 0, 0, 0, time.UTC),
			},
			want: []string{
				"2020-11-01:2020-11-01",
				"2020-11-02:2020-11-08",
				"2020-11-09:2020-11-15",
				"2020-11-16:2020-11-22",
				"2020-11-23:2020-11-29",
				"2020-11-30:2020-12-06",
				"2020-12-07:2020-12-11",
			},
		},
		{
			name: "it should detect 6 weeks between 2 months, ignoring weekend only periods",
			period: Period{
				Start: time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 12, 11, 0, 0, 0, 0, time.UTC),
			},
			optFuncs: []PeriodOptionsFunc{
				IgnoreWeekendOnlyPeriods(true),
			},
			want: []string{
				"2020-11-02:2020-11-08",
				"2020-11-09:2020-11-15",
				"2020-11-16:2020-11-22",
				"2020-11-23:2020-11-29",
				"2020-11-30:2020-12-06",
				"2020-12-07:2020-12-11",
			},
		},
		{
			name: "it should detect 6 weeks between 2 months, starting on sunday",
			period: Period{
				Start: time.Date(2020, 11, 1, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 12, 11, 0, 0, 0, 0, time.UTC),
			},
			optFuncs: []PeriodOptionsFunc{
				WeekStartsOnSunday(true),
			},
			want: []string{
				"2020-11-01:2020-11-07",
				"2020-11-08:2020-11-14",
				"2020-11-15:2020-11-21",
				"2020-11-22:2020-11-28",
				"2020-11-29:2020-12-05",
				"2020-12-06:2020-12-11",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := WeeksOnPeriod(tc.period, tc.optFuncs...)
			// convert dates to string to easily compare/diff it
			gotStr := make([]string, len(got))
			for i, t := range got {
				gotStr[i] = fmt.Sprintf("%s:%s", t.Start.Format("2006-01-02"), t.End.Format("2006-01-02"))
			}
			if !reflect.DeepEqual(gotStr, tc.want) {
				t.Errorf(diff.Cmp(tc.want, gotStr))
			}
		})
	}
}

func TestMonthsOnPeriod(t *testing.T) {
	cases := []struct {
		name     string
		period   Period
		optFuncs []PeriodOptionsFunc
		want     []string
	}{
		{
			name: "it should detect 3 months",
			period: Period{
				Start: time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 12, 11, 0, 0, 0, 0, time.UTC),
			},
			want: []string{
				"2020-10-13:2020-10-31",
				"2020-11-01:2020-11-30",
				"2020-12-01:2020-12-11",
			},
		},
		{
			name: "it should detect 2 months, ignoring weekend only periods",
			period: Period{
				Start: time.Date(2020, 10, 31, 0, 0, 0, 0, time.UTC),
				End:   time.Date(2020, 12, 11, 0, 0, 0, 0, time.UTC),
			},
			optFuncs: []PeriodOptionsFunc{
				IgnoreWeekendOnlyPeriods(true),
			},
			want: []string{
				"2020-11-01:2020-11-30",
				"2020-12-01:2020-12-11",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := MonthsOnPeriod(tc.period, tc.optFuncs...)
			// convert dates to string to easily compare/diff it
			gotStr := make([]string, len(got))
			for i, t := range got {
				gotStr[i] = fmt.Sprintf("%s:%s", t.Start.Format("2006-01-02"), t.End.Format("2006-01-02"))
			}
			if !reflect.DeepEqual(gotStr, tc.want) {
				t.Errorf(diff.Cmp(tc.want, gotStr))
			}
		})
	}
}
