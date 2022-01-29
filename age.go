package person

import (
	"errors"
	"time"
)

// TODO: define format
// TODO: AgeOn and Age should return error on invalid format

var errDobIsAfterDate = errors.New("age on: date of birth is after provided date")

// Age returns persons age formatted using format.
// For example 31 years, 2 months, 1 week, and 2 days.
func Age(dob time.Time, format string) string {
	return ""
}

// AgeOn returns persons age on a specific date formatted using format.
// It returns an error when provided date is before the date of birth (dob).
// For example 31 years, 2 months, 1 week, and 2 days.
func AgeOn(dob, date time.Time, format string) (string, error) {
	if dob.After(date) {
		return "", errDobIsAfterDate
	}
	return "2 days", nil
}

func IsAdult(dob time.Time, years time.Duration) bool {
	return false
}
