package person_test

import (
	"testing"
	"time"

	"github.com/antklim/person"
)

// TODO: move test data to CSV (especially TestCalculateDateDiff data)
// TODO: use CSV test data in all tests that calculate and output age

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
	testCases := []struct {
		dob    time.Time
		format string
		err    string
	}{
		{
			dob:    time.Now().Add(time.Second),
			format: "%D",
			err:    "date of birth is in the future",
		},
		{
			dob:    time.Now().AddDate(-2, 1, 0),
			format: " %Z m",
			err:    `format " %Z m" has unknown verb Z`,
		},
	}

	for _, tC := range testCases {
		got, err := person.Age(tC.dob, tC.format)
		if err == nil {
			t.Fatalf("Age(%s, %s) = %s, want to fail due to %s", tC.dob.Format(dateFmt),
				tC.format, got, tC.err)
		} else {
			if err.Error() != tC.err {
				t.Errorf("Age(%s, %s) failed: %v, want to fail due to %s",
					tC.dob.Format(dateFmt), tC.format, err, tC.err)
			}
			if got != "" {
				t.Errorf("Age(%s, %s) = %s, want formatted date to be an empty string",
					tC.dob.Format(dateFmt), tC.format, got)
			}
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
	dob := time.Date(1991, time.April, 1, 13, 17, 0, 0, time.UTC)

	testCases := []struct {
		date   time.Time
		format string
		err    string
	}{
		{
			date:   dob.Add(-time.Second),
			format: "%D",
			err:    "date of birth is in the future",
		},
		{
			date:   dob.AddDate(1, 1, 1),
			format: " %G %f_+",
			err:    `format " %G %f_+" has unknown verb G`,
		},
	}

	for _, tC := range testCases {
		got, err := person.AgeOn(dob, tC.date, tC.format)
		if err == nil {
			t.Fatalf("AgeOn(%s, %s, %s) = %s, want to fail due to %s", dob.Format(dateFmt),
				tC.date.Format(dateFmt), tC.format, got, tC.err)
		} else {
			if err.Error() != tC.err {
				t.Errorf("AgeOn(%s, %s, %s) failed: %v, want to fail due to %s",
					dob.Format(dateFmt), tC.date.Format(dateFmt), tC.format, err, tC.err)
			}
			if got != "" {
				t.Errorf("AgeOn(%s, %s, %s) = %s, want formatted date to be an empty string",
					dob.Format(dateFmt), tC.date.Format(dateFmt), tC.format, got)
			}
		}
	}
}

// TODO: add tests with %y, %m, %w, %d and custom time units names

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

func TestAgeFormatParse(t *testing.T) {
	testCases := []struct {
		format   string
		expected person.DateDiffFormat
	}{
		{
			format:   "   %Y   ",
			expected: person.DateDiffFormat{HasYear: true},
		},
		{
			format:   "   %y   ",
			expected: person.DateDiffFormat{HasYear: true, YearValueOnly: true},
		},
		{
			format:   "%y    %Y", // if verb repeated the latest value will be used
			expected: person.DateDiffFormat{HasYear: true},
		},
		{
			format:   "   %M   ",
			expected: person.DateDiffFormat{HasMonth: true},
		},
		{
			format:   "   %m   ",
			expected: person.DateDiffFormat{HasMonth: true, MonthValueOnly: true},
		},
		{
			format:   "%y    %Y", // if verb repeated the latest value will be used
			expected: person.DateDiffFormat{HasYear: true},
		},
		{
			format:   "   %W   ",
			expected: person.DateDiffFormat{HasWeek: true},
		},
		{
			format:   "   %w   ",
			expected: person.DateDiffFormat{HasWeek: true, WeekValueOnly: true},
		},
		{
			format:   "%y    %Y", // if verb repeated the latest value will be used
			expected: person.DateDiffFormat{HasYear: true},
		},
		{
			format:   "   %D   ",
			expected: person.DateDiffFormat{HasDay: true},
		},
		{
			format:   "   %d   ",
			expected: person.DateDiffFormat{HasDay: true, DayValueOnly: true},
		},
		{
			format:   "%d    %D", // if verb repeated the latest value will be used
			expected: person.DateDiffFormat{HasDay: true},
		},
		{
			format: "%Y  %m%D",
			expected: person.DateDiffFormat{
				HasYear:  true,
				HasMonth: true, MonthValueOnly: true,
				HasDay: true,
			},
		},
		{
			format: "  %Y%W%d",
			expected: person.DateDiffFormat{
				HasYear: true,
				HasWeek: true,
				HasDay:  true, DayValueOnly: true,
			},
		},
		{
			format: " %y%M%w %D ",
			expected: person.DateDiffFormat{
				HasYear: true, YearValueOnly: true,
				HasMonth: true,
				HasWeek:  true, WeekValueOnly: true,
				HasDay: true,
			},
		},
		{
			format:   "  %y%d  ",
			expected: person.DateDiffFormat{HasYear: true, YearValueOnly: true, HasDay: true, DayValueOnly: true},
		},
		{
			format:   "X",
			expected: person.DateDiffFormat{},
		},
	}
	for _, tC := range testCases {
		got, err := person.UnmarshalDateDiffFormat(tC.format)
		if err != nil {
			t.Errorf("UnmarshalAgeFormat(%s) failed: %v", tC.format, err)
		} else if got != tC.expected {
			t.Errorf("UnmarshalAgeFormat(%s) = %v, want %v", tC.format, got, tC.expected)
		}
	}
}

func TestAgeFormatParseFails(t *testing.T) {
	format := "%X%L %S"
	expected := `format "%X%L %S" has unknown verb X`

	got, err := person.UnmarshalDateDiffFormat(format)
	if err == nil {
		t.Fatalf("UnmarshalAgeFormat(%s) = %v, want to fail due to %s", format, got, expected)
	} else if err.Error() != expected {
		t.Errorf("UnmarshalAgeFormat(%s) failed: %v, want to fail due to %s", format, err, expected)
	}
}

// nolint:funlen
func TestCalculateDateDiff(t *testing.T) {
	// years months
	// years weeks
	// years days
	// months weeks
	// months days
	// weeks days
	// years months weeks
	// years months days
	// years weeks days
	// months weeks days
	// years months weeks days

	baseDate := time.Date(2000, time.April, 17, 0, 0, 0, 0, time.UTC)
	testCases := []struct {
		start    time.Time
		end      time.Time
		format   person.DateDiffFormat
		expected person.DateDiff
	}{
		// 2000-04-17 - 2003-03-16
		// SINGLE UNITS
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasYear: true},
			expected: person.DateDiff{Years: 2},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasMonth: true},
			expected: person.DateDiff{Months: 34},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasWeek: true},
			expected: person.DateDiff{Weeks: 151},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasDay: true},
			expected: person.DateDiff{Days: 1063},
		},

		// UNITS DOUBLES
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasYear: true, HasMonth: true},
			expected: person.DateDiff{Years: 2, Months: 10},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasYear: true, HasWeek: true},
			expected: person.DateDiff{Years: 2, Weeks: 47},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasYear: true, HasDay: true},
			expected: person.DateDiff{Years: 2, Days: 333},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasMonth: true, HasWeek: true},
			expected: person.DateDiff{Months: 34, Weeks: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasMonth: true, HasDay: true},
			expected: person.DateDiff{Months: 34, Days: 27},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasWeek: true, HasDay: true},
			expected: person.DateDiff{Weeks: 151, Days: 6},
		},

		// UNIT TRIPLETS
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasYear: true, HasMonth: true, HasWeek: true},
			expected: person.DateDiff{Years: 2, Months: 10, Weeks: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasYear: true, HasMonth: true, HasDay: true},
			expected: person.DateDiff{Years: 2, Months: 10, Days: 27},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasYear: true, HasWeek: true, HasDay: true},
			expected: person.DateDiff{Years: 2, Weeks: 47, Days: 4},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasMonth: true, HasWeek: true, HasDay: true},
			expected: person.DateDiff{Months: 34, Weeks: 3, Days: 6},
		},

		// UNIT QUARTERS
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   person.DateDiffFormat{HasYear: true, HasMonth: true, HasWeek: true, HasDay: true},
			expected: person.DateDiff{Years: 2, Months: 10, Weeks: 3, Days: 6},
		},

		// 2000-04-17, 2003-04-17
		// SINGLE UNIT
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasYear: true},
			expected: person.DateDiff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasMonth: true},
			expected: person.DateDiff{Months: 36},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasWeek: true},
			expected: person.DateDiff{Weeks: 156},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasDay: true},
			expected: person.DateDiff{Days: 1095},
		},

		// UNITS DOUBLES
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasYear: true, HasMonth: true},
			expected: person.DateDiff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasYear: true, HasWeek: true},
			expected: person.DateDiff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasYear: true, HasDay: true},
			expected: person.DateDiff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasMonth: true, HasWeek: true},
			expected: person.DateDiff{Months: 36},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasMonth: true, HasDay: true},
			expected: person.DateDiff{Months: 36},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasWeek: true, HasDay: true},
			expected: person.DateDiff{Weeks: 156, Days: 3},
		},

		// UNIT TRIPLETS
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasYear: true, HasMonth: true, HasWeek: true},
			expected: person.DateDiff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasYear: true, HasMonth: true, HasDay: true},
			expected: person.DateDiff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasYear: true, HasWeek: true, HasDay: true},
			expected: person.DateDiff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasMonth: true, HasWeek: true, HasDay: true},
			expected: person.DateDiff{Months: 36},
		},

		// UNIT QUARTERS
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   person.DateDiffFormat{HasYear: true, HasMonth: true, HasWeek: true, HasDay: true},
			expected: person.DateDiff{Years: 3},
		},
	}
	for _, tC := range testCases {
		got := person.CalculateDateDiff(tC.start, tC.end, tC.format)
		if got != tC.expected {
			t.Errorf("CalculateDateDiff(%s, %s, %v) = %v, want %v",
				tC.start.Format(dateFmt), tC.end.Format(dateFmt), tC.format, got, tC.expected)
		}
	}
}
