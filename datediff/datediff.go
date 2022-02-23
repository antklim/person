package datediff

// TODO: add documentation comments

import (
	"fmt"
	"time"
)

const (
	monthsInYear = 12
	daysInWeek   = 7
)

// TODO: define behavior in case when rawFormat is an empty string. Options are:
//	- return an error
//	- return an empty string response and don't calculate age
//	- declare define format and return age accordingly

// %Y, %y for years
// %M, %m for months
// %W, %w for weeks
// %D, %d for days
//
// %Y, %M, %W, and %D = 5 years, 4 months, 3 weeks, and 2 days
// %y years and %w weeks = 5 years and 3 weeks
// Y years and w weeks = Y years and w weeks

// TODO: refactoring. dateDiffFormat can be replaced with the bit's mask
type Format struct {
	HasYear        bool
	YearValueOnly  bool
	HasMonth       bool
	MonthValueOnly bool
	HasWeek        bool
	WeekValueOnly  bool
	HasDay         bool
	DayValueOnly   bool
}

func Unmarshal(rawFormat string) (Format, error) {
	result := Format{}
	end := len(rawFormat)
	for i := 0; i < end; {
		for i < end && rawFormat[i] != '%' {
			i++
		}
		if i >= end {
			// done processing format string
			break
		}
		// process verb
		i++
		switch c := rawFormat[i]; c {
		case 'Y':
			result.HasYear = true
			result.YearValueOnly = false
		case 'y':
			result.HasYear = true
			result.YearValueOnly = true
		case 'M':
			result.HasMonth = true
			result.MonthValueOnly = false
		case 'm':
			result.HasMonth = true
			result.MonthValueOnly = true
		case 'W':
			result.HasWeek = true
			result.WeekValueOnly = false
		case 'w':
			result.HasWeek = true
			result.WeekValueOnly = true
		case 'D':
			result.HasDay = true
			result.DayValueOnly = false
		case 'd':
			result.HasDay = true
			result.DayValueOnly = true
		default:
			return Format{}, fmt.Errorf("format %q has unknown verb %c", rawFormat, c)
		}
	}

	return result, nil
}

type Diff struct {
	Years  int
	Months int
	Weeks  int
	Days   int
	// f      Format
}

// TODO: add months and weeks output
func (d Diff) Format(f Format) string {
	switch {
	case f.HasDay:
		return formatNoun(d.Days, "day")
	case f.HasYear:
		return formatNoun(d.Years, "year")
	default:
		return ""
	}
}

// TODO: embedd format to Diff, so that it can just call method String()
func NewDiff(start, end time.Time, f Format) Diff {
	diff := Diff{}

	if f.HasYear {
		diff.Years = fullYearsDiff(start, end)
		start = start.AddDate(diff.Years, 0, 0)
	}

	if f.HasMonth {
		// getting to the closest year to the end date to reduce
		// amount of the interations during the full month calculation
		var years int
		if !f.HasYear {
			years = fullYearsDiff(start, end)
		}
		months := fullMonthsDiff(start.AddDate(years, 0, 0), end)
		diff.Months = years*monthsInYear + months
		start = start.AddDate(0, diff.Months, 0)
	}

	if f.HasWeek {
		diff.Weeks = fullWeeksDiff(start, end)
		start = start.AddDate(0, 0, diff.Weeks*daysInWeek)
	}

	if f.HasDay {
		diff.Days = fullDaysDiff(start, end)
	}

	return diff
}

func fullYearsDiff(start, end time.Time) (years int) {
	years = end.Year() - start.Year()
	if start.AddDate(years, 0, 0).After(end) {
		years--
	}
	return
}

func fullMonthsDiff(start, end time.Time) (months int) {
	for start.AddDate(0, months+1, 0).Before(end) ||
		start.AddDate(0, months+1, 0).Equal(end) {
		months++
	}
	return
}

func fullWeeksDiff(start, end time.Time) (weeks int) {
	days := daysInWeek
	for start.AddDate(0, 0, days).Before(end) ||
		start.AddDate(0, 0, days).Equal(end) {
		weeks++
		days += daysInWeek
	}
	return
}

func fullDaysDiff(start, end time.Time) (days int) {
	for start.AddDate(0, 0, days+1).Before(end) ||
		start.AddDate(0, 0, days+1).Equal(end) {
		days++
	}
	return
}

// formatNoun takes a number n and noun s in singular form.
// It returns a number and correct form of noun (singular or plural).
func formatNoun(n int, s string) string {
	f := "%d %s"
	if n%10 != 1 || n%100 == 11 {
		f += "s"
	}
	return fmt.Sprintf(f, n, s)
}
