package person

import (
	"errors"
	"fmt"
	"time"

	"github.com/antklim/person/datediff"
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

const (
	monthsInYear = 12
	daysInWeek   = 7
)

var (
	errDobIsInTheFuture = errors.New("date of birth is in the future")
)

// var timeUnits = map[rune]string{
// 	'Y': "year",
// 	'M': "month",
// 	'W': "week",
// 	'D': "day",
// }

// Age returns persons age formatted using format. It calculates age based on
// provided date of birth (dob) and current date. It returns an error when the
// provided date of birth is in the future.
// For example 31 years, 2 months, 1 week, and 2 days.
func Age(dob time.Time, rawFormat string) (string, error) {
	format, err := datediff.Unmarshal(rawFormat)
	if err != nil {
		return "", err
	}
	now := time.Now()
	return ageOn(dob, now, format)
}

// AgeOn returns persons age on a specific date formatted using format.
// It returns an error when provided date is before the date of birth (dob).
// For example 31 years, 2 months, 1 week, and 2 days.
func AgeOn(dob, date time.Time, rawFormat string) (string, error) {
	format, err := datediff.Unmarshal(rawFormat)
	if err != nil {
		return "", err
	}
	return ageOn(dob, date, format)
}

// // IsAdult returns if a person with provided date of birth is adult.
// // adultAge in years
// // WARN: not implemented
// func IsAdult(dob time.Time, adultAge int) bool {
// 	now := time.Now()
// 	return isAdultOn(dob, now, adultAge)
// }

// // WARN: not implemented
// func IsAdultOn(dob, date time.Time, adultAge int) bool {
// 	return isAdultOn(dob, date, adultAge)
// }

func ageOn(dob, date time.Time, format datediff.Format) (string, error) {
	if dob.After(date) {
		return "", errDobIsInTheFuture
	}

	d := datediff.NewDiff(dob, date, format)
	age := d.Format(format)

	return age, nil
}

// func isAdultOn(dob, date time.Time, adultAge int) bool {
// 	panic("not implemented")
// }

// TODO: refactoring. dateDiffFormat can be replaced with the bit's mask
type dateDiffFormat struct {
	HasYear        bool
	YearValueOnly  bool
	HasMonth       bool
	MonthValueOnly bool
	HasWeek        bool
	WeekValueOnly  bool
	HasDay         bool
	DayValueOnly   bool
}

func unmarshalDateDiffFormat(rawFormat string) (dateDiffFormat, error) {
	result := dateDiffFormat{}
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
			return dateDiffFormat{}, fmt.Errorf("format %q has unknown verb %c", rawFormat, c)
		}
	}

	return result, nil
}

type dateDiff struct {
	Years  int
	Months int
	Weeks  int
	Days   int
	// f      ageFormat
}

// TODO: add months, weeks output
func (d dateDiff) Format(f dateDiffFormat) string {
	switch {
	case f.HasDay:
		return formatNoun(d.Days, "day")
	case f.HasYear:
		return formatNoun(d.Years, "year")
	default:
		return ""
	}
}

func calculateDateDiff(start, end time.Time, f dateDiffFormat) dateDiff {
	diff := dateDiff{}

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
