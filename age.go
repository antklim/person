package person

import (
	"errors"
	"time"
)

// TODO: define format
// TODO: AgeOn and Age should return error on invalid format

var (
	errDobIsInTheFuture       = errors.New("age: date of birth is in the future")
	errDobIsInTheFutureOfDate = errors.New("age on: date of birth is in the future of provided date")
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
	if dob.After(date) {
		return "", errDobIsInTheFutureOfDate
	}
	return "2 days", nil
}

func IsAdult(dob time.Time, years time.Duration) bool {
	return false
}
