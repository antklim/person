package name_test

import (
	"testing"

	"github.com/antklim/person/name"
)

var nameParts = [][]string{
	nil,
	{},
	{"", "     ", "	"},
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
