package datediff

// TODO: add documentation comments

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	monthsInYear = 12
	daysInWeek   = 7
)

var formatUnits = map[string]string{
	"%Y": "year",
	"%y": "year",
	"%M": "month",
	"%m": "months",
	"%W": "week",
	"%w": "week",
	"%D": "day",
	"%d": "day",
}

var errStartIsAfterEnd = errors.New("start date is after end date")

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
type format struct {
	HasYear        bool
	YearValueOnly  bool
	HasMonth       bool
	MonthValueOnly bool
	HasWeek        bool
	WeekValueOnly  bool
	HasDay         bool
	DayValueOnly   bool
}

func unmarshal(rawFormat string) (format, error) {
	result := format{}
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
			return format{}, fmt.Errorf("format %q has unknown verb %c", rawFormat, c)
		}
	}

	return result, nil
}

type Diff struct {
	Years     int
	Months    int
	Weeks     int
	Days      int
	rawFormat string
}

// NewDiff creates Diff according to the provided format.
func NewDiff(start, end time.Time, rawFormat string) (Diff, error) {
	if start.After(end) {
		return Diff{}, errStartIsAfterEnd
	}

	diff := Diff{rawFormat: rawFormat}

	f, err := unmarshal(rawFormat)
	if err != nil {
		return Diff{}, err
	}

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

	return diff, nil
}

func (d Diff) Equal(other Diff) bool {
	return d.Years == other.Years &&
		d.Months == other.Months &&
		d.Weeks == other.Weeks &&
		d.Days == other.Days
}

func (d Diff) Format(rawFormat string) string {
	return d.format(rawFormat)
}

func (d Diff) String() string {
	return d.format(d.rawFormat)
}

func (d Diff) format(rawFormat string) string {
	result := rawFormat

	// TODO: properly format lower case verbs %y, %m,...

	for verb, unit := range formatUnits {
		if strings.Contains(result, verb) {
			var n int
			switch unit {
			case "year":
				n = d.Years
			case "month":
				n = d.Months
			case "week":
				n = d.Weeks
			case "day":
				n = d.Days
			}
			result = strings.ReplaceAll(result, verb, formatNoun(n, unit))
		}
	}

	return result
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

// formatNoun takes a positive number n and noun s in singular form.
// It returns a number and correct form of noun (singular or plural).
func formatNoun(n int, s string) string {
	f := "%d %s"
	if n != 1 {
		f += "s"
	}
	return fmt.Sprintf(f, n, s)
}
