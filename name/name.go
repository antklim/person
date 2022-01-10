package name

import "strings"

// FullName returns a person's full name. If a join of all name parts result to
// an empty string then fallback returned. Otherwise name parts returned joined
// with a whitespace.
// Note: non printable chanracters are not handled and not trimmed.
func FullName(fallback string, parts ...string) string {
	newParts := nonEmptyParts(parts)
	fullName := strings.Join(newParts, " ")

	if fullName == "" {
		return fallback
	}

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
