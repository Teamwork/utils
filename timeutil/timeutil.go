package timeutil // import "github.com/teamwork/utils/timeutil"

import (
	"time"

	"github.com/jinzhu/now"
)

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

// Period stores the start and end dates of a period.
type Period struct {
	Start time.Time
	End   time.Time
}

// PeriodOptions store options for period related functions.
type PeriodOptions struct {
	ignoreWeekendOnly bool
	startsOnSunday    bool
	flexiblePeriod    bool
}

// PeriodOptionsFunc function to modify the period options.
type PeriodOptionsFunc func(p *PeriodOptions)

// IgnoreWeekendOnlyPeriods ignore when the period only has weekdays on a
// weekend.
func IgnoreWeekendOnlyPeriods(ignore bool) PeriodOptionsFunc {
	return PeriodOptionsFunc(func(p *PeriodOptions) {
		p.ignoreWeekendOnly = ignore
	})
}

// WeekStartsOnSunday use Sunday as beginning of the week when calculating the
// periods.
func WeekStartsOnSunday(startsOnSunday bool) PeriodOptionsFunc {
	return PeriodOptionsFunc(func(p *PeriodOptions) {
		p.startsOnSunday = startsOnSunday
	})
}

// FlexiblePeriod change the start/end date to fit the desired grouping.
func FlexiblePeriod(flexible bool) PeriodOptionsFunc {
	return PeriodOptionsFunc(func(p *PeriodOptions) {
		p.flexiblePeriod = flexible
	})
}

// WeeksOnPeriod extracts all weeks on the given period.
func WeeksOnPeriod(period Period, optFuncs ...PeriodOptionsFunc) []Period {
	return groupPeriod(
		period,
		func(referenceDate *now.Now) time.Time {
			return referenceDate.BeginningOfWeek()
		},
		func(referenceDate *now.Now) time.Time {
			return referenceDate.EndOfWeek()
		},
		func(referenceDate *now.Now) time.Time {
			return referenceDate.AddDate(0, 0, 7)
		},
		optFuncs...,
	)
}

// MonthsOnPeriod extracts all months on the given period.
func MonthsOnPeriod(period Period, optFuncs ...PeriodOptionsFunc) []Period {
	return groupPeriod(
		period,
		func(referenceDate *now.Now) time.Time {
			return referenceDate.BeginningOfMonth()
		},
		func(referenceDate *now.Now) time.Time {
			return referenceDate.EndOfMonth()
		},
		func(referenceDate *now.Now) time.Time {
			return referenceDate.AddDate(0, 1, 0)
		},
		optFuncs...,
	)
}

func groupPeriod(period Period, toStart, toEnd, add func(*now.Now) time.Time, optFuncs ...PeriodOptionsFunc) []Period {
	options := PeriodOptions{
		ignoreWeekendOnly: false,
	}
	for _, optFunc := range optFuncs {
		optFunc(&options)
	}

	referenceDate := now.With(period.Start)
	if options.startsOnSunday {
		referenceDate.WeekStartDay = time.Sunday
	} else {
		referenceDate.WeekStartDay = time.Monday
	}
	referenceDate.Time = toStart(referenceDate)

	var periods []Period

	for referenceDate.Before(period.End) {
		p := Period{
			Start: toStart(referenceDate),
			End:   toEnd(referenceDate),
		}
		if !options.flexiblePeriod {
			if p.Start.Before(period.Start) {
				p.Start = period.Start
			}
			if p.End.After(period.End) {
				p.End = period.End
			}
		}
		periods = append(periods, p)

		// move to the next period
		referenceDate.Time = add(referenceDate)
	}

	if options.ignoreWeekendOnly {
		i := 0
		for _, period := range periods {
			weekendOnly := true
			dateFunc := DateRange(period.Start, period.End)
			for d := dateFunc(); !d.IsZero(); d = dateFunc() {
				if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
					weekendOnly = false
					break
				}
			}
			if !weekendOnly {
				periods[i] = period
				i++
			}
		}
		periods = periods[:i]
	}

	return periods
}

// DateRange iterates on each day between start and end.
//
//     dateFunc := timeutil.DateRange(startDate, endDate)
//     for d := dateFunc(); !d.IsZero(); d = dateFunc() {
//       println(d.Format(time.RFC3339))
//     }
func DateRange(start, end time.Time) func() time.Time {
	// safety rounding
	start = now.With(start).BeginningOfDay()
	end = now.With(end).BeginningOfDay()

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}
