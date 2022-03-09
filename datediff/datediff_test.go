package datediff_test

import (
	"encoding/csv"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/antklim/person/datediff"
)

const dateFmt = "2006-01-02"

// these are the fields of testdata/dates.csv
const (
	startFld = iota
	endFld
	formatFld
	yearsFld
	monthsFld
	weeksFld
	daysFld
	printFld
)

type datesRecord struct {
	start  time.Time
	end    time.Time
	format string
	diff   datediff.Diff
	print  string
}

func loadDatesRecord(r []string) (datesRecord, error) {
	start, err := time.Parse(dateFmt, r[startFld])
	if err != nil {
		return datesRecord{}, err
	}
	end, err := time.Parse(dateFmt, r[endFld])
	if err != nil {
		return datesRecord{}, err
	}
	years, err := strconv.Atoi(r[yearsFld])
	if err != nil {
		return datesRecord{}, err
	}
	months, err := strconv.Atoi(r[monthsFld])
	if err != nil {
		return datesRecord{}, err
	}
	weeks, err := strconv.Atoi(r[weeksFld])
	if err != nil {
		return datesRecord{}, err
	}
	days, err := strconv.Atoi(r[daysFld])
	if err != nil {
		return datesRecord{}, err
	}

	return datesRecord{
		start:  start,
		end:    end,
		format: r[formatFld],
		diff:   datediff.Diff{Years: years, Months: months, Weeks: weeks, Days: days},
		print:  r[printFld],
	}, nil
}

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

func TestNewDiff(t *testing.T) {
	f, err := os.Open("testdata/dates.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comment = '#'
	testCases, err := r.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	for _, tC := range testCases {
		testRecord, err := loadDatesRecord(tC)
		if err != nil {
			t.Fatal(err)
		}
		got, err := datediff.NewDiff(testRecord.start, testRecord.end, testRecord.format)
		if err != nil {
			t.Errorf("NewDiff(%s, %s, %s) failed: %v",
				testRecord.start.Format(dateFmt), testRecord.end.Format(dateFmt),
				testRecord.format, err)
		} else if !got.Equal(testRecord.diff) {
			t.Errorf("NewDiff(%s, %s, %s) = %v, want %#v",
				testRecord.start.Format(dateFmt), testRecord.end.Format(dateFmt),
				testRecord.format, got, testRecord.diff)
		}
	}
}

func TestNewDiffFails(t *testing.T) {
	testCases := []struct {
		start  time.Time
		end    time.Time
		format string
		err    string
	}{
		{
			start:  time.Now().Add(time.Hour),
			end:    time.Now(),
			format: "%D",
			err:    "start date is after end date",
		},
		{
			start:  time.Now(),
			end:    time.Now().Add(time.Hour),
			format: " %Z m",
			err:    `format " %Z m" has unknown verb Z`,
		},
	}

	for _, tC := range testCases {
		got, err := datediff.NewDiff(tC.start, tC.end, tC.format)
		if err == nil {
			t.Fatalf("NewDiff(%s, %s, %s) = %v, want to fail due to %s",
				tC.start.Format(dateFmt), tC.end.Format(dateFmt), tC.format, got, tC.err)
		} else if err.Error() != tC.err {
			t.Errorf("NewDiff(%s, %s, %s) failed: %v, want to fail due to %s",
				tC.start.Format(dateFmt), tC.end.Format(dateFmt), tC.format, err, tC.err)
		}
	}
}

// TODO: add tests with %y, %m, %w, %d and custom time units names
// TODO: re-use the dates from the CSV/previous test. In this case
// coorectness of diff calculation validated.

func TestString(t *testing.T) {
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
		diff, err := datediff.NewDiff(start, end, tC.format)
		if err != nil {
			t.Errorf("NewDiff(%s, %s, %s) failed: %v",
				start.Format(dateFmt), end.Format(dateFmt), tC.format, err)
		}
		got := diff.String()
		if got != tC.expected {
			t.Errorf("String() = %s, want %s", got, tC.expected)
		}
	}
}

func TestFormat(t *testing.T) {
	start := time.Date(2000, time.April, 17, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(3, -1, -1)
	format := "%Y, %M, %W and %D"
	diff, err := datediff.NewDiff(start, end, format)
	if err != nil {
		t.Errorf("NewDiff(%s, %s, %s) failed: %v",
			start.Format(dateFmt), end.Format(dateFmt), format, err)
	}
	format = "%Y"
	got := diff.Format(format)
	expected := "2 years"
	if got != expected {
		t.Errorf("Format(%s) = %s, want %s", format, got, expected)
	}
}
