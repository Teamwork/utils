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

// StartOfMonth returns the first day of the month of date.
func StartOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

// EndOfMonth returns the last day of the month of date.
func EndOfMonth(date time.Time) time.Time {
	// go to the next month, then a day of 0 removes a day leaving us
	// at the last day of dates month.
	return time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, date.Location())
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
