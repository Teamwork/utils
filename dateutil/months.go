package dateutil

import "time"

// MonthsTo returns the number of months from the current date until "a"  Will
// return 1 if less than a full month away
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
