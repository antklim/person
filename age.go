package person

// TODO: review documentation comments

import (
	"errors"
	"time"

	"github.com/antklim/person/datediff"
)

var (
	errDobIsInTheFuture = errors.New("date of birth is in the future")
)

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
	return d.Format(format), nil
}

// func isAdultOn(dob, date time.Time, adultAge int) bool {
// 	panic("not implemented")
// }
