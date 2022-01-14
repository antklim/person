package name

import "strings"

// FullName returns a person's full name as a result of joining parts of the
// name.
// Every name part trimmed from leading and trailing white spaces. Name part
// that have value of empty string is omitted from joining. Name parts joined
// with a single white space separator.
func FullName(parts []string) string {
	p := formatAndFilter(parts, strings.TrimSpace, defaultFilter)
	fullName := strings.Join(p, " ")
	return fullName
}

// FullNameDefault returns a person's full name as a result of joining parts of
// the name. If the joining of name parts produces an empty string then default
// value d returned.
func FullNameDefault(parts []string, d string) string {
	if n := FullName(parts); n != "" {
		return n
	}
	return d
}

// FullNameFormatFunc returns a person's full name as a result of joining parts
// of the name. Every name part is formatted with format function f. Formatted
// name parts then joined.
func FullNameFormatFunc(parts []string, f func(string) string) string {
	p := formatAndFilter(parts, f, defaultFilter)
	fullName := strings.Join(p, " ")
	return fullName
}

// FullNameDefaultFormatFunc returns a person's full name as a result of joining
// parts of the name. Every name part is formatted with format function f.
// Formatted name parts then joined. If the joining of name parts produces an
// empty string then default value d returned.
func FullNameDefaultFormatFunc(parts []string, d string, f func(string) string) string {
	p := formatAndFilter(parts, f, defaultFilter)
	fullName := strings.Join(p, " ")
	if fullName == "" {
		return d
	}
	return fullName
}

type formatFunc func(string) string
type filterFunc func(string) bool

func formatAndFilter(a []string, f formatFunc, ff filterFunc) []string {
	var result []string
	for _, s := range a {
		if s := f(s); ff(s) {
			result = append(result, s)
		}
	}
	return result
}

func defaultFilter(s string) bool {
	return s != ""
}
