# person

The `person` package provides methods that allow you to work with personal information.

# Installation
`go get github.com/antklim/person`

# Usage
```go
package main

import (
	"fmt"
	"time"

	"github.com/antklim/person"
)

func main() {
	nameParts := []string{" John", "     ", "   ", "	Smith  ", "Doe"}
	fullName := person.FullName(nameParts)
	fmt.Printf("Full name: %s\n", fullName)

	dob := time.Now().AddDate(-20, 0, 0)
	age, _ := person.Age(dob, "%Y")
	fmt.Printf("Age: %s\n", age)

	adultAge := 18
	isAdult, _ := person.IsAdult(dob, adultAge)
	fmt.Printf("Is adult: %t\n", isAdult)

	// Output:
	// Full name: John Smith Doe
	// Age: 20 years
	// Is adult: true
}
```
