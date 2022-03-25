package person_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/antklim/person"
)

func ExampleFullName() {
	nameParts := []string{" Johann", "     ", "   ", " Sebastian  ", "Bach"}
	fullName := person.FullName(nameParts)
	fmt.Println(fullName)

	// Output:
	// Johann Sebastian Bach
}

func ExampleFullNameDefault() {
	nameParts := []string{" Johann", "     ", "   ", " Sebastian  ", "Bach"}
	fullName := person.FullNameDefault(nameParts, "unknown")
	fmt.Println(fullName)

	nameParts = []string{"", "     ", " "}
	fullName = person.FullNameDefault(nameParts, "unknown")
	fmt.Println(fullName)

	// Output:
	// Johann Sebastian Bach
	// unknown
}

func ExampleFullNameFormatFunc() {
	f := func(s string) string {
		s = strings.TrimSpace(s)
		if s == "Sebastian" {
			s = "-"
		}
		return s
	}

	nameParts := []string{" Johann", "     ", "   ", " Sebastian  ", "Bach"}
	fullName := person.FullNameFormatFunc(nameParts, f)
	fmt.Println(fullName)

	// Output:
	// Johann - Bach
}

func ExampleAgeOn() {
	dob, _ := time.Parse("2006-01-02", "2000-01-01")
	ondate, _ := time.Parse("2006-01-02", "2003-03-16")
	age, _ := person.AgeOn(dob, ondate, "%Y %M %D")
	fmt.Println(age)

	// Output:
	// 3 years 2 months 15 days
}

func ExampleIsAdult() {
	dob, _ := time.Parse("2006-01-02", "2000-01-01")
	adultAge := 18
	isAdult, _ := person.IsAdult(dob, adultAge)
	fmt.Println(isAdult)

	// Output:
	// true
}
