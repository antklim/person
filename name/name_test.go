package name_test

import (
	"strings"
	"testing"

	"github.com/antklim/person/name"
)

var nameParts = [][]string{
	nil,
	{},
	{"", "     ", " "},
	{" Johann", "     ", "   ", "	Sebastian  ", "Bach"},
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
		got := name.FullName(v)
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
		got := name.FullNameDefault(v, dflt)
		if got != want {
			t.Errorf("FullNameDefault(%v, %s) = %s, want %s", v, dflt, got, want)
		}
	}
}

func TestFullNameFormatFunc(t *testing.T) {
	f := func(s string) string {
		if s == " " {
			return "John"
		}
		return strings.TrimSpace(s)
	}

	expected := []string{
		"",
		"",
		"John",
		"Johann Sebastian Bach",
	}
	for i, v := range nameParts {
		want := expected[i]
		got := name.FullNameFormatFunc(v, f)
		if got != want {
			t.Errorf("FullNameFormatFunc(%v, formatFunc) = %s, want %s", v, got, want)
		}
	}
}
