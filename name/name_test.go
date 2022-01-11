package name_test

import (
	"testing"

	"github.com/antklim/person/name"
)

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
