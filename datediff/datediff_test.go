package datediff_test

import (
	"testing"
	"time"

	"github.com/antklim/person/datediff"
)

const dateFmt = "2006-01-02"

func TestUnmarshal(t *testing.T) {
	testCases := []struct {
		format   string
		expected datediff.Format
	}{
		{
			format:   "   %Y   ",
			expected: datediff.Format{HasYear: true},
		},
		{
			format:   "   %y   ",
			expected: datediff.Format{HasYear: true, YearValueOnly: true},
		},
		{
			format:   "%y    %Y", // if verb repeated the latest value will be used
			expected: datediff.Format{HasYear: true},
		},
		{
			format:   "   %M   ",
			expected: datediff.Format{HasMonth: true},
		},
		{
			format:   "   %m   ",
			expected: datediff.Format{HasMonth: true, MonthValueOnly: true},
		},
		{
			format:   "%y    %Y", // if verb repeated the latest value will be used
			expected: datediff.Format{HasYear: true},
		},
		{
			format:   "   %W   ",
			expected: datediff.Format{HasWeek: true},
		},
		{
			format:   "   %w   ",
			expected: datediff.Format{HasWeek: true, WeekValueOnly: true},
		},
		{
			format:   "%y    %Y", // if verb repeated the latest value will be used
			expected: datediff.Format{HasYear: true},
		},
		{
			format:   "   %D   ",
			expected: datediff.Format{HasDay: true},
		},
		{
			format:   "   %d   ",
			expected: datediff.Format{HasDay: true, DayValueOnly: true},
		},
		{
			format:   "%d    %D", // if verb repeated the latest value will be used
			expected: datediff.Format{HasDay: true},
		},
		{
			format: "%Y  %m%D",
			expected: datediff.Format{
				HasYear:  true,
				HasMonth: true, MonthValueOnly: true,
				HasDay: true,
			},
		},
		{
			format: "  %Y%W%d",
			expected: datediff.Format{
				HasYear: true,
				HasWeek: true,
				HasDay:  true, DayValueOnly: true,
			},
		},
		{
			format: " %y%M%w %D ",
			expected: datediff.Format{
				HasYear: true, YearValueOnly: true,
				HasMonth: true,
				HasWeek:  true, WeekValueOnly: true,
				HasDay: true,
			},
		},
		{
			format:   "  %y%d  ",
			expected: datediff.Format{HasYear: true, YearValueOnly: true, HasDay: true, DayValueOnly: true},
		},
		{
			format:   "X",
			expected: datediff.Format{},
		},
	}
	for _, tC := range testCases {
		got, err := datediff.Unmarshal(tC.format)
		if err != nil {
			t.Errorf("Unmarshal(%s) failed: %v", tC.format, err)
		} else if got != tC.expected {
			t.Errorf("Unmarshal(%s) = %v, want %v", tC.format, got, tC.expected)
		}
	}
}

func TestUnmarshalFails(t *testing.T) {
	format := "%X%L %S"
	expected := `format "%X%L %S" has unknown verb X`

	got, err := datediff.Unmarshal(format)
	if err == nil {
		t.Fatalf("Unmarshal(%s) = %v, want to fail due to %s", format, got, expected)
	} else if err.Error() != expected {
		t.Errorf("Unmarshal(%s) failed: %v, want to fail due to %s", format, err, expected)
	}
}

// TODO: move test data to CSV

func TestNewDiff(t *testing.T) { // nolint:funlen
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
		format   datediff.Format
		expected datediff.Diff
	}{
		// 2000-04-17 - 2003-03-16
		// SINGLE UNITS
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasYear: true},
			expected: datediff.Diff{Years: 2},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasMonth: true},
			expected: datediff.Diff{Months: 34},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasWeek: true},
			expected: datediff.Diff{Weeks: 151},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasDay: true},
			expected: datediff.Diff{Days: 1063},
		},

		// UNITS DOUBLES
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasYear: true, HasMonth: true},
			expected: datediff.Diff{Years: 2, Months: 10},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasYear: true, HasWeek: true},
			expected: datediff.Diff{Years: 2, Weeks: 47},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasYear: true, HasDay: true},
			expected: datediff.Diff{Years: 2, Days: 333},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasMonth: true, HasWeek: true},
			expected: datediff.Diff{Months: 34, Weeks: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasMonth: true, HasDay: true},
			expected: datediff.Diff{Months: 34, Days: 27},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasWeek: true, HasDay: true},
			expected: datediff.Diff{Weeks: 151, Days: 6},
		},

		// UNIT TRIPLETS
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasYear: true, HasMonth: true, HasWeek: true},
			expected: datediff.Diff{Years: 2, Months: 10, Weeks: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasYear: true, HasMonth: true, HasDay: true},
			expected: datediff.Diff{Years: 2, Months: 10, Days: 27},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasYear: true, HasWeek: true, HasDay: true},
			expected: datediff.Diff{Years: 2, Weeks: 47, Days: 4},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasMonth: true, HasWeek: true, HasDay: true},
			expected: datediff.Diff{Months: 34, Weeks: 3, Days: 6},
		},

		// UNIT QUARTERS
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, -1, -1),
			format:   datediff.Format{HasYear: true, HasMonth: true, HasWeek: true, HasDay: true},
			expected: datediff.Diff{Years: 2, Months: 10, Weeks: 3, Days: 6},
		},

		// 2000-04-17, 2003-04-17
		// SINGLE UNIT
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasYear: true},
			expected: datediff.Diff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasMonth: true},
			expected: datediff.Diff{Months: 36},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasWeek: true},
			expected: datediff.Diff{Weeks: 156},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasDay: true},
			expected: datediff.Diff{Days: 1095},
		},

		// UNITS DOUBLES
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasYear: true, HasMonth: true},
			expected: datediff.Diff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasYear: true, HasWeek: true},
			expected: datediff.Diff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasYear: true, HasDay: true},
			expected: datediff.Diff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasMonth: true, HasWeek: true},
			expected: datediff.Diff{Months: 36},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasMonth: true, HasDay: true},
			expected: datediff.Diff{Months: 36},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasWeek: true, HasDay: true},
			expected: datediff.Diff{Weeks: 156, Days: 3},
		},

		// UNIT TRIPLETS
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasYear: true, HasMonth: true, HasWeek: true},
			expected: datediff.Diff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasYear: true, HasMonth: true, HasDay: true},
			expected: datediff.Diff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasYear: true, HasWeek: true, HasDay: true},
			expected: datediff.Diff{Years: 3},
		},
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasMonth: true, HasWeek: true, HasDay: true},
			expected: datediff.Diff{Months: 36},
		},

		// UNIT QUARTERS
		{
			start:    baseDate,
			end:      baseDate.AddDate(3, 0, 0),
			format:   datediff.Format{HasYear: true, HasMonth: true, HasWeek: true, HasDay: true},
			expected: datediff.Diff{Years: 3},
		},

		{
			start:    time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
			end:      time.Date(2010, time.April, 14, 0, 0, 0, 0, time.UTC),
			format:   datediff.Format{HasYear: true, HasDay: true},
			expected: datediff.Diff{Years: 10, Days: 103},
		},
	}
	for _, tC := range testCases {
		got := datediff.NewDiff(tC.start, tC.end, tC.format)
		if got != tC.expected {
			t.Errorf("CalculateDateDiff(%s, %s, %v) = %v, want %v",
				tC.start.Format(dateFmt), tC.end.Format(dateFmt), tC.format, got, tC.expected)
		}
	}
}

// TODO: add tests with %y, %m, %w, %d and custom time units names
// TODO: re-use the dates from the CSV/previous test. In this case
// coorectness of diff calculation validated.

func TestFormatPrint(t *testing.T) {
	start := time.Date(2000, time.April, 17, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(3, -1, -1)
	testCases := []struct {
		format   string
		expected string
	}{
		{
			format:   "%Y",
			expected: "2 years",
		},
		{
			format:   "%M",
			expected: "34 months",
		},
		{
			format:   "%W",
			expected: "151 weeks",
		},
		{
			format:   "%D",
			expected: "1063 days",
		},
		{
			format:   "%Y and %M",
			expected: "2 years and 10 months",
		},
		{
			format:   "%Y and %W",
			expected: "2 years and 47 weeks",
		},
		{
			format:   "%Y and %D",
			expected: "2 years and 333 days",
		},
		{
			format:   "%M and %W",
			expected: "34 months and 3 weeks",
		},
		{
			format:   "%M and %D",
			expected: "34 months and 27 days",
		},
		{
			format:   "%W and %D",
			expected: "151 weeks and 6 days",
		},
		{
			format:   "%Y, %M and %W",
			expected: "2 years, 10 months and 3 weeks",
		},
		{
			format:   "%Y, %M and %D",
			expected: "2 years, 10 months and 27 days",
		},
		{
			format:   "%Y, %W and %D",
			expected: "2 years, 47 weeks and 4 days",
		},
		{
			format:   "%M, %W and %D",
			expected: "34 months, 3 weeks and 6 days",
		},
		{
			format:   "%Y, %M, %W and %D",
			expected: "2 years, 10 months, 3 weeks and 6 days",
		},
	}
	for _, tC := range testCases {
		format, err := datediff.Unmarshal(tC.format)
		if err != nil {
			t.Errorf("Unmarshal(%s) failed: %v", tC.format, err)
		}

		diff := datediff.NewDiff(start, end, format)
		got := diff.Format(format, tC.format)
		if got != tC.expected {
			t.Errorf("Format(%#v, %s) = %s, want %s", format, tC.format, got, tC.expected)
		}
	}
}