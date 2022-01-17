package person_test

import (
	"strings"
	"testing"

	"github.com/antklim/person"
)

var nameParts = [][]string{
	nil,
	{},
	{"", "     ", " "},
	{" Johann", "     ", "   ", " Sebastian  ", "Bach"},
}

func formatFunc(s string) string {
	if s == " " {
		return "John"
	}
	return strings.TrimSpace(s)
}

func TestFullName(t *testing.T) {
	expected := []string{
		"",
		"",
		"",
		"Johann Sebastian Bach",
	}

	for i, v := range nameParts {
		want := expected[i]
		got := person.FullName(v)
		if got != want {
			t.Errorf("FullName(%v) = %s, want %s", v, got, want)
		}
	}
}

func TestFullNameDefault(t *testing.T) {
	dflt := "John Doe"
	expected := []string{
		dflt,
		dflt,
		dflt,
		"Johann Sebastian Bach",
	}

	for i, v := range nameParts {
		want := expected[i]
		got := person.FullNameDefault(v, dflt)
		if got != want {
			t.Errorf("FullNameDefault(%v, %s) = %s, want %s", v, dflt, got, want)
		}
	}
}

func TestFullNameFormatFunc(t *testing.T) {
	expected := []string{
		"",
		"",
		"John",
		"Johann Sebastian Bach",
	}
	for i, v := range nameParts {
		want := expected[i]
		got := person.FullNameFormatFunc(v, formatFunc)
		if got != want {
			t.Errorf("FullNameFormatFunc(%v, formatFunc) = %s, want %s", v, got, want)
		}
	}
}

func TestFullNameDefaultFormatFunc(t *testing.T) {
	dflt := "John Doe"
	expected := []string{
		dflt,
		dflt,
		"John",
		"Johann Sebastian Bach",
	}
	for i, v := range nameParts {
		want := expected[i]
		got := person.FullNameDefaultFormatFunc(v, dflt, formatFunc)
		if got != want {
			t.Errorf("FullNameDefaultFormatFunc(%v, %s, formatFunc) = %s, want %s", v, dflt, got, want)
		}
	}
}
