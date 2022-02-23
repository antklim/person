package person

import (
	"errors"
	"fmt"
	"time"
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
	hoursInDay   = 24
	hoursInYear  = 365 * hoursInDay
	monthsInYear = 12
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
	format, err := unmarshalAgeFormat(rawFormat)
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
	format, err := unmarshalAgeFormat(rawFormat)
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

func ageOn(dob, date time.Time, f ageFormat) (string, error) {
	if dob.After(date) {
		return "", errDobIsInTheFuture
	}

	d := date.Sub(dob)
	age := formatDuration(d, f)

	return age, nil
}

// func isAdultOn(dob, date time.Time, adultAge int) bool {
// 	panic("not implemented")
// }

func formatDuration(d time.Duration, f ageFormat) string {
	// this function is responsible for printing the result
	// result calculation should be done in another function - such as parse duration

	// it should know the state of output
	// for example, when format %Y %M %D - it outputs full years, months and days
	// but when format is just %D - it outputs age in days (3653 days - 10 years and 3 days)

	switch {
	case f.HasDay:
		days := int(d.Hours() / hoursInDay)
		return formatNoun(days, "day")
	case f.HasYear:
		years := int(d.Hours() / hoursInYear)
		return formatNoun(years, "year")
	default:
		return ""
	}
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

// TODO: refactoring. ageFormat can be replaced with the bit's mask
type ageFormat struct {
	HasYear        bool
	YearValueOnly  bool
	HasMonth       bool
	MonthValueOnly bool
	HasWeek        bool
	WeekValueOnly  bool
	HasDay         bool
	DayValueOnly   bool
}

func unmarshalAgeFormat(format string) (ageFormat, error) {
	result := ageFormat{}
	end := len(format)
	for i := 0; i < end; {
		for i < end && format[i] != '%' {
			i++
		}
		if i >= end {
			// done processing format string
			break
		}
		// process verb
		i++
		switch c := format[i]; c {
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
			return ageFormat{}, fmt.Errorf("format %q has unknown verb %c", format, c)
		}
	}

	return result, nil
}

type dateDiff struct {
	Years  int
	Months int
	Weeks  int
	Days   int
}

func calculateDateDiff(start, end time.Time, f ageFormat) dateDiff {
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
