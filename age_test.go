package person_test

import (
	"errors"
	"testing"
	"time"

	"github.com/antklim/person"
)

var testDOB = time.Date(1991, time.April, 1, 13, 17, 0, 0, time.UTC)

func TestAge(t *testing.T) {
	t.Skip("not implemented")
	// testCases := []struct {
	// 	format   string
	// 	expected string
	// }{
	// 	{
	// 		format:   "",
	// 		expected: "",
	// 	},
	// }
}

func TestAgeFails(t *testing.T) {
	format := "D"
	dob := time.Now().Add(time.Second)
	expected := errors.New("date of birth is in the future")

	got, err := person.Age(dob, format)
	if err == nil {
		t.Fatalf("Age(%s, %s) = %s, want to fail due to %v", dob.Format(time.RFC3339),
			format, got, expected)
	} else {
		if err.Error() != expected.Error() {
			t.Errorf("Age(%s, %s) failed: %v, want to fail due to %v",
				dob.Format(time.RFC3339), format, err, expected)
		}
		if got != "" {
			t.Errorf("Age(%s, %s) = %s, want formatted date to be an empty string",
				dob.Format(time.RFC3339), format, got)
		}
	}
}

func TestAgeOn(t *testing.T) {
	testCases := []struct {
		date     time.Time
		format   string
		expected string
	}{
		// TODO: add leap year test
		{
			date:     time.Date(1991, time.April, 3, 13, 17, 0, 0, time.UTC),
			format:   "D",
			expected: "2 days",
		},
		{
			date:     time.Date(1991, time.April, 2, 13, 17, 0, 0, time.UTC),
			format:   "D",
			expected: "1 day",
		},
		{
			date:     time.Date(1992, time.April, 2, 0, 0, 0, 0, time.UTC),
			format:   "Y",
			expected: "1 year",
		},
		{
			date:     time.Date(1993, time.April, 2, 0, 0, 0, 0, time.UTC),
			format:   "Y",
			expected: "2 years",
		},
		{
			date:     time.Date(2003, time.April, 2, 0, 0, 0, 0, time.UTC),
			format:   "Y",
			expected: "12 years",
		},
		{
			date:     time.Date(2012, time.April, 2, 0, 0, 0, 0, time.UTC),
			format:   "Y",
			expected: "21 year",
		},
	}

	for _, tC := range testCases {
		got, err := person.AgeOn(testDOB, tC.date, tC.format)
		if err != nil {
			t.Errorf("AgeOn(%s, %s, %s) failed: %v", testDOB.Format(time.RFC3339),
				tC.date.Format(time.RFC3339), tC.format, err)
		} else if got != tC.expected {
			t.Errorf("AgeOn(%s, %s, %s) = %s, want %s", testDOB.Format(time.RFC3339),
				tC.date.Format(time.RFC3339), tC.format, got, tC.expected)
		}
	}
}

func TestAgeOnFails(t *testing.T) {
	format := "D"
	date := testDOB.Add(-time.Second)
	expected := errors.New("date of birth is after provided date")

	got, err := person.AgeOn(testDOB, date, format)
	if err == nil {
		t.Fatalf("AgeOn(%s, %s, %s) = %s, want to fail due to %v", testDOB.Format(time.RFC3339),
			date.Format(time.RFC3339), format, got, expected)
	} else {
		if err.Error() != expected.Error() {
			t.Errorf("AgeOn(%s, %s, %s) failed: %v, want to fail due to %v",
				testDOB.Format(time.RFC3339), date.Format(time.RFC3339), format, err, expected)
		}
		if got != "" {
			t.Errorf("AgeOn(%s, %s, %s) = %s, want formatted date to be an empty string",
				testDOB.Format(time.RFC3339), date.Format(time.RFC3339), format, got)
		}
	}
}