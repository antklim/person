package person

import (
	"errors"
	"fmt"
	"time"
)

// TODO: define format
// TODO: AgeOn and Age should return error on invalid format

// %Y, %y for years
// %M, %m for months
// %W, %w for weeks
// %D, %d for days
// %H, %h for hours
//
// %Y, %M, %W, and %D = 5 years, 4 months, 3 weeks, and 2 days
// %y years and %w weeks = 5 years and 3 weeks
// Y years and w weeks = Y years and w weeks
// %Z years = years (unknown verb replaced with '')

const (
	hoursInDay  = 24
	hoursInYear = 365 * hoursInDay
	// hoursInLeapYear = 366 * hoursInDay
)

var (
	errDobIsInTheFuture       = errors.New("date of birth is in the future")
	errDobIsInTheFutureOfDate = errors.New("date of birth is after provided date")
)

// Age returns persons age formatted using format. It calculates age based on
// provided date of birth (dob) and current date. It returns an error when the
// provided date of birth is in the future.
// For example 31 years, 2 months, 1 week, and 2 days.
func Age(dob time.Time, format string) (string, error) {
	now := time.Now()
	if dob.After(now) {
		return "", errDobIsInTheFuture
	}
	return "", nil
}

// AgeOn returns persons age on a specific date formatted using format.
// It returns an error when provided date is before the date of birth (dob).
// For example 31 years, 2 months, 1 week, and 2 days.
func AgeOn(dob, date time.Time, format string) (string, error) {
	return ageOn(dob, date, format)
}

func IsAdult(dob time.Time, years time.Duration) bool {
	return false
}

func ageOn(dob, date time.Time, format string) (string, error) {
	if dob.After(date) {
		return "", errDobIsInTheFutureOfDate
	}

	d := date.Sub(dob)
	age := formatDuration(d, format)

	return age, nil
}

func formatDuration(d time.Duration, format string) string {
	switch format {
	case "D":
		days := int(d.Hours() / hoursInDay)
		if days%10 == 1 {
			return fmt.Sprintf("%d day", days)
		}
		return fmt.Sprintf("%d days", days)
	case "Y":
		years := int(d.Hours() / hoursInYear)
		if years%10 == 1 {
			return fmt.Sprintf("%d year", years)
		}
		return fmt.Sprintf("%d years", years)
	default:
		return ""
	}
}
