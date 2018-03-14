// Package timeutil provides helpers for working with time.Time.
package timeutil // import "github.com/teamwork/utils/timeutil"

import "time"

// UnixMilli returns the number of milliseconds elapsed since January 1, 1970
// UTC.
func UnixMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

// DaysBetween return the number of whole days between a start date and end date
func DaysBetween(fromDate, toDate time.Time) int {
	return int(toDate.Sub(fromDate) / (24 * time.Hour))
}

// StartOfDay returns the start of t's day.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of t's day.
func EndOfDay(t time.Time) time.Time {
	// Finding the end of the day is tricky, due to leap seconds (We can't
	// assume that 23:59:59 is the last second--or that it even happens on any
	// given day due to negative leap seconds) and DST.
	//
	// The strategy here is to find the start of the day, then add 36 hours
	// (must be more than 24, to account for DST), then find midnight of that
	// date, and subtract one tick.
	tomorrow := StartOfDay(t).Add(36 * time.Hour)
	return StartOfDay(tomorrow).Add(-time.Nanosecond)
}

// Tomorrow returns the start of the day after t.
func Tomorrow(t time.Time) time.Time {
	return StartOfDay(StartOfDay(t).Add(36 * time.Hour))
}

// StartOfMonth returns the first day of t's month.
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the end of the last day of t's month.
func EndOfMonth(t time.Time) time.Time {
	// go to the next month, then a day of 0 removes a day leaving us
	// at the last day of dates month.
	return EndOfDay(time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()))
}

// FormatAsZulu gets a ISO 8601 formatted date. The date is assumed to be in
// UTC ("Zulu time").
//
// TODO: I think we can just use t.Format(time.RFC3339)?
func FormatAsZulu(t time.Time) string {
	return t.Format("2006-01-02T15:04:05Z")
}

// MonthsTo returns the number of months from the current date to the given
// date. The number of months is always rounded down, with a minimal value of 1.
//
// For example this returns 2:
//     MonthsTo(time.Now().Add(24 * time.Hour * 70))
//
// Dates in the past are not supported, and their behaviour is undefined!
func MonthsTo(a time.Time) int {
	var days int
	startDate := time.Now()
	lastDayOfYear := func(t time.Time) time.Time {
		return time.Date(t.Year(), 12, 31, 0, 0, 0, 0, t.Location())
	}

	firstDayOfNextYear := func(t time.Time) time.Time {
		return time.Date(t.Year()+1, 1, 1, 0, 0, 0, 0, t.Location())
	}
	cur := startDate
	for cur.Year() < a.Year() {
		// add 1 to count the last day of the year too.
		days += lastDayOfYear(cur).YearDay() - cur.YearDay() + 1
		cur = firstDayOfNextYear(cur)
	}
	days += a.YearDay() - cur.YearDay()
	if startDate.AddDate(0, 0, days).After(a) {
		days--
	}
	months := (days / 30)
	if months == 0 {
		months = 1
	}
	return months
}

// The below functions are modified from https://github.com/jinzhu/now which
// is released under the MIT license available here:
// https://opensource.org/licenses/MIT

// StartOfWeek returns the start of t's week (Monday).
func StartOfWeek(t time.Time) time.Time {
	daysToSubtract := t.Weekday() - 1
	if daysToSubtract == -1 { // Sunday
		daysToSubtract = 6
	}
	return StartOfDay(t).Add(time.Duration(-daysToSubtract) * 24 * time.Hour)
}

// EndOfWeek returns the end of t's week (Sunday).
func EndOfWeek(t time.Time) time.Time {
	if t.Weekday() == time.Sunday {
		return EndOfDay(t)
	}

	daysToAdd := 7 - t.Weekday()
	return EndOfDay(t).Add(time.Duration(daysToAdd) * 24 * time.Hour)
}

// StartOfQuarter returns the first day of t's quarter.
func StartOfQuarter(t time.Time) time.Time {
	month := StartOfMonth(t)
	offset := (int(month.Month()) - 1) % 3
	return month.AddDate(0, -offset, 0)
}

// EndOfQuarter returns the end of the last day of t's quarter.
func EndOfQuarter(t time.Time) time.Time {
	return StartOfQuarter(t).AddDate(0, 3, 0).Add(-time.Nanosecond)
}
