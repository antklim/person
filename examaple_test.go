package person_test

import (
	"fmt"
	"strings"

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
