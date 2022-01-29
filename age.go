package person

import "time"

// Age returns persons age formatted using format.
// For example 31 years, 2 months, 1 week, and 2 days.
func Age(dob time.Time, format string) string {
	return ""
}

func AgeIn(dob, in time.Time, format string) (string, error) {
	return "", nil
}

func IsAdult(dob time.Time, years time.Duration) bool {
	return false
}
