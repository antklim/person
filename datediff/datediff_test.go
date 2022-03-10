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

type datediffRecord struct {
	start  time.Time
	end    time.Time
	format string
	diff   datediff.Diff
	print  string
}

func loadDatediffRecord(r []string) (datediffRecord, error) {
	start, err := time.Parse(dateFmt, r[startFld])
	if err != nil {
		return datediffRecord{}, err
	}
	end, err := time.Parse(dateFmt, r[endFld])
	if err != nil {
		return datediffRecord{}, err
	}
	years, err := strconv.Atoi(r[yearsFld])
	if err != nil {
		return datediffRecord{}, err
	}
	months, err := strconv.Atoi(r[monthsFld])
	if err != nil {
		return datediffRecord{}, err
	}
	weeks, err := strconv.Atoi(r[weeksFld])
	if err != nil {
		return datediffRecord{}, err
	}
	days, err := strconv.Atoi(r[daysFld])
	if err != nil {
		return datediffRecord{}, err
	}

	return datediffRecord{
		start:  start,
		end:    end,
		format: r[formatFld],
		diff:   datediff.Diff{Years: years, Months: months, Weeks: weeks, Days: days},
		print:  r[printFld],
	}, nil
}

func loadDatediffRecordsForTest() ([]datediffRecord, error) {
	f, err := os.Open("testdata/datediff.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comment = '#'
	rawRecords, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var dr []datediffRecord
	for _, rr := range rawRecords {
		r, err := loadDatediffRecord(rr)
		if err != nil {
			return nil, err
		}
		dr = append(dr, r)
	}
	return dr, nil
}

func TestUnmarshal(t *testing.T) {
	testCases := []struct {
		format   string
		expected datediff.Format
	}{
		{
			format: "   %Y   ",
			expected: datediff.Format{
				UnitsMask: datediff.HasYearMask,
			},
		},
		{
			format: "   %y   ",
			expected: datediff.Format{
				YearValueOnly: true,
				UnitsMask:     datediff.HasYearMask,
			},
		},
		{
			format: "%y    %Y", // if verb repeated the latest value will be used
			expected: datediff.Format{
				UnitsMask: datediff.HasYearMask,
			},
		},
		{
			format: "   %M   ",
			expected: datediff.Format{
				HasMonth:  true,
				UnitsMask: datediff.HasMonthMask,
			},
		},
		{
			format: "   %m   ",
			expected: datediff.Format{
				HasMonth:       true,
				MonthValueOnly: true,
				UnitsMask:      datediff.HasMonthMask,
			},
		},
		{
			format: "%y    %Y", // if verb repeated the latest value will be used
			expected: datediff.Format{
				UnitsMask: datediff.HasYearMask,
			},
		},
		{
			format: "   %W   ",
			expected: datediff.Format{
				HasWeek:   true,
				UnitsMask: datediff.HasWeekMask,
			},
		},
		{
			format: "   %w   ",
			expected: datediff.Format{
				HasWeek:       true,
				WeekValueOnly: true,
				UnitsMask:     datediff.HasWeekMask,
			},
		},
		{
			format: "%y    %Y", // if verb repeated the latest value will be used
			expected: datediff.Format{
				UnitsMask: datediff.HasYearMask,
			},
		},
		{
			format: "   %D   ",
			expected: datediff.Format{
				HasDay:    true,
				UnitsMask: datediff.HasDayMask,
			},
		},
		{
			format: "   %d   ",
			expected: datediff.Format{
				HasDay:       true,
				DayValueOnly: true,
				UnitsMask:    datediff.HasDayMask,
			},
		},
		{
			format: "%d    %D", // if verb repeated the latest value will be used
			expected: datediff.Format{
				HasDay:    true,
				UnitsMask: datediff.HasDayMask,
			},
		},
		{
			format: "%Y  %m%D",
			expected: datediff.Format{
				HasMonth: true, MonthValueOnly: true,
				HasDay:    true,
				UnitsMask: datediff.HasYearMask | datediff.HasMonthMask | datediff.HasDayMask,
			},
		},
		{
			format: "  %Y%W%d",
			expected: datediff.Format{
				HasWeek: true,
				HasDay:  true, DayValueOnly: true,
				UnitsMask: datediff.HasYearMask | datediff.HasWeekMask | datediff.HasDayMask,
			},
		},
		{
			format: " %y%M%w %D ",
			expected: datediff.Format{
				YearValueOnly: true,
				HasMonth:      true,
				HasWeek:       true, WeekValueOnly: true,
				HasDay:    true,
				UnitsMask: datediff.HasYearMask | datediff.HasMonthMask | datediff.HasWeekMask | datediff.HasDayMask,
			},
		},
		{
			format: "  %y%d  ",
			expected: datediff.Format{
				YearValueOnly: true,
				HasDay:        true,
				DayValueOnly:  true,
				UnitsMask:     datediff.HasYearMask | datediff.HasDayMask,
			},
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
	testCases, err := loadDatediffRecordsForTest()
	if err != nil {
		t.Fatal(err)
	}

	for _, tC := range testCases {
		got, err := datediff.NewDiff(tC.start, tC.end, tC.format)
		if err != nil {
			t.Errorf("NewDiff(%s, %s, %s) failed: %v",
				tC.start.Format(dateFmt), tC.end.Format(dateFmt), tC.format, err)
		} else if !got.Equal(tC.diff) {
			t.Errorf("NewDiff(%s, %s, %s) = %v, want %#v",
				tC.start.Format(dateFmt), tC.end.Format(dateFmt), tC.format, got, tC.diff)
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

func TestString(t *testing.T) {
	testCases, err := loadDatediffRecordsForTest()
	if err != nil {
		t.Fatal(err)
	}

	for _, tC := range testCases {
		diff, err := datediff.NewDiff(tC.start, tC.end, tC.format)
		if err != nil {
			t.Errorf("NewDiff(%s, %s, %s) failed: %v",
				tC.start.Format(dateFmt), tC.end.Format(dateFmt), tC.format, err)
		}
		got := diff.String()
		if got != tC.print {
			t.Errorf("String() = %s, want %s", got, tC.print)
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
