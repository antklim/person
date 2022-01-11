package name

import "strings"

// FullName returns a person's full name as a result of joining parts of the
// name.
// Every name part trimmed from leading and trailing white spaces. Name part
// that have value of empty string is omitted from joining. Name parts joined
// with a single white space separator.
func FullName(parts []string) string {
	newParts := nonEmptyParts(parts)
	fullName := strings.Join(newParts, " ")
	return fullName
}

func nonEmptyPart(p string) bool {
	return p != ""
}

func nonEmptyParts(parts []string) []string {
	var filteredParts []string
	for _, p := range parts {
		if p := strings.TrimSpace(p); nonEmptyPart(p) {
			filteredParts = append(filteredParts, p)
		}
	}
	return filteredParts
}

// FullNameDefault returns a person's full name as a result of joining parts of
// the name. If the joining of name parts produces an empty string then default
// value d returned.
func FullNameDefault(parts []string, d string) string {
	return ""
}

// FullNameFormatFunc returns a person's full name as a result of joining parts
// of the name. Every name part is formatted with format function f. Formatted
// name parts then joined.
func FullNameFormatFunc(parts []string, f func(string) string) string {
	return ""
}

// FullNameDefaultFormatFunc returns a person's full name as a result of joining
// parts of the name. Every name part is formatted with format function f.
// Formatted name parts then joined. If the joining of name parts produces an
// empty string then default value d returned.
func FullNameDefaultFormatFunc(parts []string, d string, f func(string) string) string {
	return ""
}
