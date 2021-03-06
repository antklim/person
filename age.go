package person

import (
	"errors"
	"time"

	"github.com/antklim/datediff"
)

var (
	errDobIsInTheFuture = errors.New("date of birth is in the future")
)

// Age returns persons age formatted using format. It calculates age based on
// provided date of birth (dob) and current date. It returns an error when the
// provided date of birth is in the future.
// For example 31 years, 2 months, 1 week, and 2 days.
func Age(dob time.Time, rawFormat string) (string, error) {
	now := time.Now()
	return ageOn(dob, now, rawFormat)
}

// AgeOn returns persons age on a specific date formatted using format.
// It returns an error when provided date is before the date of birth (dob).
// For example 31 years, 2 months, 1 week, and 2 days.
func AgeOn(dob, date time.Time, rawFormat string) (string, error) {
	return ageOn(dob, date, rawFormat)
}

// IsAdult returns if a person with provided date of birth is adult.
// Adult age parameter is the minimum amount of full year a
// person should have to be considered an adult. If the person's age in years
// is greater or equal to the adult age the function returns true.
// It returns an error when the provided date of birth is in the future.
func IsAdult(dob time.Time, adultAge int) (bool, error) {
	return isAdultOn(dob, time.Now(), adultAge)
}

// IsAdultOn returns if a person with provided date of birth is adult on a
// specific date. Adult age parameter is the minimum amount of full year a
// person should have to be considered an adult. If the person's age in years
// is greater or equal to the adult age the function returns true.
// It returns an error when provided date is before the date of birth (dob).
func IsAdultOn(dob, date time.Time, adultAge int) (bool, error) {
	return isAdultOn(dob, date, adultAge)
}

func ageOn(dob, date time.Time, rawFormat string) (string, error) {
	if dob.After(date) {
		return "", errDobIsInTheFuture
	}
	age, err := datediff.NewDiff(dob, date, rawFormat)
	if err != nil {
		return "", err
	}
	return age.String(), nil
}

func isAdultOn(dob, date time.Time, adultAge int) (bool, error) {
	if dob.After(date) {
		return false, errDobIsInTheFuture
	}
	age, err := datediff.NewDiff(dob, time.Now(), "%y")
	if err != nil {
		return false, err
	}
	return age.Years >= adultAge, nil
}
