package person_test

import (
	"encoding/csv"
	"os"
	"testing"
	"time"

	"github.com/antklim/person"
)

const dateFmt = "2006-01-02"

// these are the fields of testdata/age.csv
const (
	dobFld = iota
	ondateFld
	formatFld
	printFld
)

type ageRecord struct {
	dob    time.Time
	ondate time.Time
	format string
	print  string
}

func loadAgeRecord(r []string) (ageRecord, error) {
	dob, err := time.Parse(dateFmt, r[dobFld])
	if err != nil {
		return ageRecord{}, err
	}
	ondate, err := time.Parse(dateFmt, r[ondateFld])
	if err != nil {
		return ageRecord{}, err
	}

	return ageRecord{
		dob:    dob,
		ondate: ondate,
		format: r[formatFld],
		print:  r[printFld],
	}, nil
}

func loadAgeRecordsForTest() ([]ageRecord, error) {
	f, err := os.Open("testdata/age.csv")
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

	var dr []ageRecord
	for _, rr := range rawRecords {
		r, err := loadAgeRecord(rr)
		if err != nil {
			return nil, err
		}
		dr = append(dr, r)
	}
	return dr, nil
}

var testInvalidAgeInputs = []struct {
	dob      time.Time
	ondate   time.Time
	format   string
	expected string
}{
	{
		dob:      time.Now().Add(24 * time.Hour),
		ondate:   time.Now(),
		format:   "%D",
		expected: "date of birth is in the future",
	},
	{
		format:   " %Z m",
		expected: `format " %Z m" has unknown verb Z`,
	},
}

func TestAge(t *testing.T) {
	now := time.Now()
	dob := now.AddDate(-11, 1, 2)
	format := "%Y"

	want := "10 years"
	got, err := person.Age(dob, format)
	if err != nil {
		t.Errorf("Age(%s, %s) failed: %v", dob.Format(dateFmt), format, err)
	} else if got != want {
		t.Logf("now: %s", now.Format(dateFmt))
		t.Errorf("Age(%s, %s) = %s, want %s", dob.Format(dateFmt), format, got, want)
	}
}

func TestAgeFails(t *testing.T) {
	for _, tC := range testInvalidAgeInputs {
		got, err := person.Age(tC.dob, tC.format)
		if err == nil {
			t.Fatalf("Age(%s, %s) = %s, want to fail due to %s", tC.dob.Format(dateFmt),
				tC.format, got, tC.expected)
		} else {
			if err.Error() != tC.expected {
				t.Errorf("Age(%s, %s) failed: %v, want to fail due to %s",
					tC.dob.Format(dateFmt), tC.format, err, tC.expected)
			}
			if got != "" {
				t.Errorf("Age(%s, %s) = %s, want formatted date to be an empty string",
					tC.dob.Format(dateFmt), tC.format, got)
			}
		}
	}
}

func TestAgeOn(t *testing.T) {
	testCases, err := loadAgeRecordsForTest()
	if err != nil {
		t.Fatal(err)
	}

	for _, tC := range testCases {
		got, err := person.AgeOn(tC.dob, tC.ondate, tC.format)
		if err != nil {
			t.Errorf("AgeOn(%s, %s, %s) failed: %v", tC.dob.Format(dateFmt),
				tC.ondate.Format(dateFmt), tC.format, err)
		} else if got != tC.print {
			t.Errorf("AgeOn(%s, %s, %s) = %s, want %s", tC.dob.Format(dateFmt),
				tC.ondate.Format(dateFmt), tC.format, got, tC.print)
		}
	}
}

func TestAgeOnFails(t *testing.T) {
	for _, tC := range testInvalidAgeInputs {
		got, err := person.AgeOn(tC.dob, tC.ondate, tC.format)
		if err == nil {
			t.Fatalf("AgeOn(%s, %s, %s) = %s, want to fail due to %s", tC.dob.Format(dateFmt),
				tC.ondate.Format(dateFmt), tC.format, got, tC.expected)
		} else {
			if err.Error() != tC.expected {
				t.Errorf("AgeOn(%s, %s, %s) failed: %v, want to fail due to %s",
					tC.dob.Format(dateFmt), tC.ondate.Format(dateFmt), tC.format, err, tC.expected)
			}
			if got != "" {
				t.Errorf("AgeOn(%s, %s, %s) = %s, want formatted date to be an empty string",
					tC.dob.Format(dateFmt), tC.ondate.Format(dateFmt), tC.format, got)
			}
		}
	}
}

func TestIsAdult(t *testing.T) {
	dob := time.Now().AddDate(-20, 0, 0)

	testCases := []struct {
		adultAge int
		expected bool
	}{
		{
			adultAge: 18,
			expected: true,
		},
		{
			adultAge: 21,
			expected: false,
		},
	}
	for _, tC := range testCases {
		got, err := person.IsAdult(dob, tC.adultAge)
		if err != nil {
			t.Errorf("IsAdult(%s, %d) failed: %v", dob.Format(dateFmt), tC.adultAge, err)
		} else if got != tC.expected {
			t.Errorf("IsAdult(%s, %d) = %t, want %t", dob.Format(dateFmt),
				tC.adultAge, got, tC.expected)
		}
	}
}

func TestIsAdultFails(t *testing.T) {
	dob := time.Now().Add(time.Hour)
	adultAge := 18
	expected := "date of birth is in the future"
	got, err := person.IsAdult(dob, adultAge)
	if err == nil {
		t.Fatalf("IsAdult(%s, %d) = %t, want to fail due %s",
			dob.Format(dateFmt), adultAge, got, expected)
	} else {
		if err.Error() != expected {
			t.Errorf("IsAdult(%s, %d) failed: %v, want to fail due to %s",
				dob.Format(dateFmt), adultAge, err, expected)
		}
		if got != false {
			t.Errorf("IsAdult(%s, %d) = %t, want to be false",
				dob.Format(dateFmt), adultAge, got)
		}
	}
}

func TestIsAdultOn(t *testing.T) {
	dob := time.Now().AddDate(-20, 0, 0)
	ondate := time.Now().AddDate(-2, 0, 0)

	testCases := []struct {
		adultAge int
		expected bool
	}{
		{
			adultAge: 18,
			expected: true,
		},
		{
			adultAge: 21,
			expected: false,
		},
	}
	for _, tC := range testCases {
		got, err := person.IsAdultOn(dob, ondate, tC.adultAge)
		if err != nil {
			t.Errorf("IsAdultOn(%s, %s, %d) failed: %v",
				dob.Format(dateFmt), ondate.Format(dateFmt), tC.adultAge, err)
		} else if got != tC.expected {
			t.Errorf("IsAdultOn(%s, %s, %d) = %t, want %t",
				dob.Format(dateFmt), ondate.Format(dateFmt), tC.adultAge, got, tC.expected)
		}
	}
}

func TestIsAdultOnFails(t *testing.T) {
	dob := time.Now()
	ondate := time.Now().AddDate(0, 0, -1)
	adultAge := 18
	expected := "date of birth is in the future"
	got, err := person.IsAdultOn(dob, ondate, adultAge)
	if err == nil {
		t.Fatalf("IsAdultOn(%s, %s, %d) = %t, want to fail due %s",
			dob.Format(dateFmt), ondate.Format(dateFmt), adultAge, got, expected)
	} else {
		if err.Error() != expected {
			t.Errorf("IsAdultOn(%s, %s, %d) failed: %v, want to fail due to %s",
				dob.Format(dateFmt), ondate.Format(dateFmt), adultAge, err, expected)
		}
		if got != false {
			t.Errorf("IsAdultOn(%s, %s, %d) = %t, want to be false",
				dob.Format(dateFmt), ondate.Format(dateFmt), adultAge, got)
		}
	}
}
