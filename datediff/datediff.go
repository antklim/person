package datediff

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
//	- declare default format and return age accordingly

// %Y, %y for years
// %M, %m for months
// %W, %w for weeks
// %D, %d for days

type format struct {
	DiffMode   DiffMode
	FormatMode FormatMode
}

type DiffMode uint8

const (
	ModeYear DiffMode = 1 << iota
	ModeMonth
	ModeWeek
	ModeDay
)

type FormatMode uint8

const (
	FormatYearValueOnly FormatMode = 1 << iota
	FormatMonthValueOnly
	FormatWeekValueOnly
	FormatDayValueOnly
)

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
			result.FormatMode &^= FormatYearValueOnly
			result.DiffMode |= ModeYear
		case 'y':
			result.FormatMode |= FormatYearValueOnly
			result.DiffMode |= ModeYear
		case 'M':
			result.FormatMode &^= FormatMonthValueOnly
			result.DiffMode |= ModeMonth
		case 'm':
			result.FormatMode |= FormatMonthValueOnly
			result.DiffMode |= ModeMonth
		case 'W':
			result.FormatMode &^= FormatWeekValueOnly
			result.DiffMode |= ModeWeek
		case 'w':
			result.FormatMode |= FormatWeekValueOnly
			result.DiffMode |= ModeWeek
		case 'D':
			result.FormatMode &^= FormatDayValueOnly
			result.DiffMode |= ModeDay
		case 'd':
			result.FormatMode |= FormatDayValueOnly
			result.DiffMode |= ModeDay
		default:
			return format{}, fmt.Errorf("format %q has unknown verb %c", rawFormat, c)
		}
	}

	return result, nil
}

// Diff describes dates difference in years, months, weeks, and days.
type Diff struct {
	Years     int
	Months    int
	Weeks     int
	Days      int
	rawFormat string
}

// NewDiff creates Diff according to the provided format.
// Provided format should contain special "verbs" that define dates difference
// calculattion logic. These are supported format verbs:
//	%Y - to calculate dates difference in years
//	%M - to calculate dates difference in months
//	%W - to calculate dates difference in weeks
//	%D - to calculate dates difference in days
//
// When format contains multiple "verbs" the date diffrence will be calculated
// starting from longest time unit to shortest. For example:
//
//	start, _ := time.Parse("2006-01-02", "2000-04-17")
//	end, _ := time.Parse("2006-01-02", "2003-03-16")
//	diff1, _ := NewDiff(start, end, "%Y")
//	diff2, _ := NewDiff(start, end, "%M")
//	diff3, _ := NewDiff(start, end, "%Y %M")
//	fmt.Println(diff1) // 2 years
//	fmt.Println(diff2) // 34 month
//	fmt.Println(diff3) // 2 years 10 months
//
// NewDiff returns error when start date is after end date or in case of invalid
// format. Format considered invalid when it contains unsupported "verb".
func NewDiff(start, end time.Time, rawFormat string) (Diff, error) {
	if start.After(end) {
		return Diff{}, errStartIsAfterEnd
	}

	diff := Diff{rawFormat: rawFormat}

	f, err := unmarshal(rawFormat)
	if err != nil {
		return Diff{}, err
	}

	if f.DiffMode&ModeYear != 0 {
		diff.Years = fullYearsDiff(start, end)
		start = start.AddDate(diff.Years, 0, 0)
	}

	if f.DiffMode&ModeMonth != 0 {
		// getting to the closest year to the end date to reduce
		// amount of the interations during the full month calculation
		var years int
		if f.DiffMode&ModeYear == 0 {
			years = fullYearsDiff(start, end)
		}
		months := fullMonthsDiff(start.AddDate(years, 0, 0), end)
		diff.Months = years*monthsInYear + months
		start = start.AddDate(0, diff.Months, 0)
	}

	if f.DiffMode&ModeWeek != 0 {
		diff.Weeks = fullWeeksDiff(start, end)
		start = start.AddDate(0, 0, diff.Weeks*daysInWeek)
	}

	if f.DiffMode&ModeDay != 0 {
		diff.Days = fullDaysDiff(start, end)
	}

	return diff, nil
}

// Equal returns true when two dates differences are equal.
func (d Diff) Equal(other Diff) bool {
	return d.Years == other.Years &&
		d.Months == other.Months &&
		d.Weeks == other.Weeks &&
		d.Days == other.Days
}

// Format formats dates difference accordig to provided format.
func (d Diff) Format(rawFormat string) string {
	return d.format(rawFormat)
}

// String formats dates difference according to the format provided at
// initialization of dates difference.
func (d Diff) String() string {
	return d.format(d.rawFormat)
}

func (d Diff) format(rawFormat string) string {
	result := rawFormat

	// TODO: properly format lower case verbs %y, %m,...
	// TODO: add feature to trim verb when unit value is 0

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
