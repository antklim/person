package datediff

import (
	"errors"
	"fmt"
	"time"
)

const (
	monthsInYear = 12
	daysInWeek   = 7
)

var (
	errStartIsAfterEnd   = errors.New("start date is after end date")
	errUndefinedDiffMode = errors.New("undefined dates difference mode")
)

type DiffMode uint8

const (
	ModeYears DiffMode = 1 << (8 - 1 - iota)
	ModeMonths
	ModeWeeks
	ModeDays
)

func unmarshal(rawFormat string) (DiffMode, error) {
	var mode DiffMode
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
		case 'Y', 'y':
			mode |= ModeYears
		case 'M', 'm':
			mode |= ModeMonths
		case 'W', 'w':
			mode |= ModeWeeks
		case 'D', 'd':
			mode |= ModeDays
		default:
			return 0, fmt.Errorf("format %q has unknown verb %c", rawFormat, c)
		}
	}

	if mode == 0 {
		return 0, errUndefinedDiffMode
	}

	return mode, nil
}

// Diff describes dates difference in years, months, weeks, and days.
type Diff struct {
	Years     int
	Months    int
	Weeks     int
	Days      int
	rawFormat string // initial format, i.e "%Y and %M"
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
// NewDiff returns error in the following cases:
//	start date is after end date
//	format contains unsupported "verb"
//	undefined dates difference mode (it happens when the format does not contain any of the supported "verbs")
func NewDiff(start, end time.Time, rawFormat string) (Diff, error) {
	if start.After(end) {
		return Diff{}, errStartIsAfterEnd
	}

	mode, err := unmarshal(rawFormat)
	if err != nil {
		return Diff{}, err
	}

	diff := Diff{rawFormat: rawFormat}

	if mode&ModeYears != 0 {
		diff.Years = fullYearsDiff(start, end)
		start = start.AddDate(diff.Years, 0, 0)
	}

	if mode&ModeMonths != 0 {
		// getting to the closest year to the end date to reduce
		// amount of the interations during the full month calculation
		var years int
		if mode&ModeYears == 0 {
			years = fullYearsDiff(start, end)
		}
		months := fullMonthsDiff(start.AddDate(years, 0, 0), end)
		diff.Months = years*monthsInYear + months
		start = start.AddDate(0, diff.Months, 0)
	}

	if mode&ModeWeeks != 0 {
		diff.Weeks = fullWeeksDiff(start, end)
		start = start.AddDate(0, 0, diff.Weeks*daysInWeek)
	}

	if mode&ModeDays != 0 {
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
func (d Diff) Format(rawFormat string) (string, error) {
	_, err := unmarshal(rawFormat)
	if err != nil {
		return "", err
	}
	return formatWithZeros(d, rawFormat), nil
}

// String formats dates difference according to the format provided at
// initialization of dates difference.
func (d Diff) String() string {
	return formatWithZeros(d, d.rawFormat)
}
