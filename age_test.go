package person_test

import (
	"errors"
	"testing"
	"time"

	"github.com/antklim/person"
)

const dateFmt = "2006-01-02"

func TestAge(t *testing.T) {
	now := time.Now()
	format := "%Y"
	testCases := []struct {
		dob      time.Time
		expected string
	}{
		{
			dob:      now.AddDate(-2, 1, 2),
			expected: "1 year",
		},
		{
			dob:      now.AddDate(-3, 1, 2),
			expected: "2 years",
		},
		{
			dob:      now.AddDate(-11, 1, 2),
			expected: "10 years",
		},
		{
			dob:      now.AddDate(-12, 1, 2),
			expected: "11 years",
		},
		{
			dob:      now.AddDate(-13, 1, 2),
			expected: "12 years",
		},
		{
			dob:      now.AddDate(-22, 1, 2),
			expected: "21 year",
		},
		{
			dob:      now.AddDate(-112, 1, 2),
			expected: "111 years",
		},
	}
	for _, tC := range testCases {
		got, err := person.Age(tC.dob, format)
		if err != nil {
			t.Logf("now: %s", now.Format(dateFmt))
			t.Errorf("Age(%s, %s) failed: %v", tC.dob.Format(dateFmt), format, err)
		} else if got != tC.expected {
			t.Logf("now: %s", now.Format(dateFmt))
			t.Errorf("Age(%s, %s) = %s, want %s", tC.dob.Format(dateFmt), format,
				got, tC.expected)
		}
	}
}

func TestAgeFails(t *testing.T) {
	format := "%D"
	dob := time.Now().Add(time.Second)
	expected := errors.New("date of birth is in the future")

	got, err := person.Age(dob, format)
	if err == nil {
		t.Fatalf("Age(%s, %s) = %s, want to fail due to %v", dob.Format(dateFmt),
			format, got, expected)
	} else {
		if err.Error() != expected.Error() {
			t.Errorf("Age(%s, %s) failed: %v, want to fail due to %v",
				dob.Format(dateFmt), format, err, expected)
		}
		if got != "" {
			t.Errorf("Age(%s, %s) = %s, want formatted date to be an empty string",
				dob.Format(dateFmt), format, got)
		}
	}
}

func TestAgeOn(t *testing.T) {
	dob := time.Date(1991, time.April, 1, 13, 17, 0, 0, time.UTC)
	testCases := []struct {
		date     time.Time
		format   string
		expected string
	}{
		// TODO: add leap year test
		{
			date:     time.Date(1991, time.April, 3, 13, 17, 0, 0, time.UTC),
			format:   "%D",
			expected: "2 days",
		},
		{
			date:     time.Date(1991, time.April, 2, 13, 17, 0, 0, time.UTC),
			format:   "%D",
			expected: "1 day",
		},
		{
			date:     time.Date(1992, time.April, 2, 0, 0, 0, 0, time.UTC),
			format:   "%Y",
			expected: "1 year",
		},
		{
			date:     time.Date(1993, time.April, 2, 0, 0, 0, 0, time.UTC),
			format:   "%Y",
			expected: "2 years",
		},
		{
			date:     time.Date(2003, time.April, 2, 0, 0, 0, 0, time.UTC),
			format:   "%Y",
			expected: "12 years",
		},
		{
			date:     time.Date(2012, time.April, 2, 0, 0, 0, 0, time.UTC),
			format:   "%Y",
			expected: "21 year",
		},
	}

	for _, tC := range testCases {
		got, err := person.AgeOn(dob, tC.date, tC.format)
		if err != nil {
			t.Errorf("AgeOn(%s, %s, %s) failed: %v", dob.Format(dateFmt),
				tC.date.Format(dateFmt), tC.format, err)
		} else if got != tC.expected {
			t.Errorf("AgeOn(%s, %s, %s) = %s, want %s", dob.Format(dateFmt),
				tC.date.Format(dateFmt), tC.format, got, tC.expected)
		}
	}
}

func TestAgeOnFails(t *testing.T) {
	format := "%D"
	dob := time.Date(1991, time.April, 1, 13, 17, 0, 0, time.UTC)
	date := dob.Add(-time.Second)
	expected := errors.New("date of birth is in the future")

	got, err := person.AgeOn(dob, date, format)
	if err == nil {
		t.Fatalf("AgeOn(%s, %s, %s) = %s, want to fail due to %v", dob.Format(dateFmt),
			date.Format(dateFmt), format, got, expected)
	} else {
		if err.Error() != expected.Error() {
			t.Errorf("AgeOn(%s, %s, %s) failed: %v, want to fail due to %v",
				dob.Format(dateFmt), date.Format(dateFmt), format, err, expected)
		}
		if got != "" {
			t.Errorf("AgeOn(%s, %s, %s) = %s, want formatted date to be an empty string",
				dob.Format(dateFmt), date.Format(dateFmt), format, got)
		}
	}
}

// TODO: add tests with months and weeks
// TODO: add tests with %y, %m, %w, %d, %h and custom time units names

func TestFormatPrint(t *testing.T) {
	t.Skip()
	dob := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	date := time.Date(2010, time.April, 14, 12, 0, 0, 0, time.UTC)
	testCases := []struct {
		format   string
		expected string
	}{
		{
			format:   "%Y",
			expected: "10 years",
		},
		{
			format:   "%D",
			expected: "3756 days",
		},
		{
			format:   "%H",
			expected: "90156 hours",
		},
		{
			format:   "%Y and %D",
			expected: "10 years and 103 days",
		},
		{
			format:   "%Y and %H",
			expected: "10 years and 2484 hours",
		},
		{
			format:   "%Y, %D, and %H",
			expected: "10 years, 103 days, and 12 hours",
		},
		{
			format:   "%D and %H",
			expected: "3756 days and 12 hours",
		},
	}
	for _, tC := range testCases {
		got, err := person.AgeOn(dob, date, tC.format)
		if err != nil {
			t.Errorf("AgeOn(%s, %s, %s) failed: %v", dob.Format(dateFmt),
				date.Format(dateFmt), tC.format, err)
		} else if got != tC.expected {
			t.Errorf("AgeOn(%s, %s, %s) = %s, want %s", dob.Format(dateFmt),
				date.Format(dateFmt), tC.format, got, tC.expected)
		}
	}
}

// TODO: add month, week, hour formatting

func TestFormatParse(t *testing.T) {
	// testCases := []struct {
	// 	format   string
	// 	expected AgeFormat
	// }{
	// 	{
	// 		format:   "%Y",
	// 		expected: AgeFormat{year: true, dyear: true},
	// 	},
	// 	{
	// 		format:   "%y",
	// 		expected: AgeFormat{year: true, dyear: false},
	// 	},
	// 	{
	// 		format:   "%D",
	// 		expected: AgeFormat{day: true, dday: true},
	// 	},
	// 	{
	// 		format:   "%d",
	// 		expected: AgeFormat{day: true, dday: false},
	// 	},
	// }
	// for _, tC := range testCases {
	// 	t.Run(tC.desc, func(t *testing.T) {

	// 	})
	// }
}
