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
	testCases := []struct {
		desc  string
		parts []string
		want  string
	}{
		{
			desc: "returns an empty string when no name parts provided",
		},
		{
			desc: "returns an empty string when all name parts trimmed to empty strings",
			parts: []string{"", "     ", "	"},
		},
		{
			desc: "returns joined name parts",
			parts: []string{" Johann", "     ", "   ", "	Sebastian  ", "Bach"},
			want: "Johann Sebastian Bach",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := name.FullName(tC.parts)
			if got != tC.want {
				t.Errorf("FullName(%v) = %s, want %s", tC.parts, got, tC.want)
			}
		})
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
