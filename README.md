# person

The `person` package provides methods that allow you to work with personal information.

# Installation
`go get github.com/antklim/person`

# Usage
```go
package main

import(
  "fmt"
  
  "github.com/antklim/person/name"
)

func main() {
  nameParts := []string{" Johann", "     ", "   ", "	Sebastian  ", "Bach"}
  fullName := name.FullName(nameParts)
  fmt.Print(fullName)
  // Output:
  // Johann Sebastian Bach
}
```
