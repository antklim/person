package name_test

import (
	"testing"

	"github.com/antklim/person/name"
)

func TestFullName(t *testing.T) {
	fallback := "John Doe"
	testCases := []struct {
		desc  string
		parts []string
		want  string
	}{
		{
			desc: "returns fallback value when no name parts provided",
			want: fallback,
		},
		{
			desc: "returns fallback value when all name parts trimmed to empty strings",
			parts: []string{"", "     ", "	"},
			want: fallback,
		},
		{
			desc: "returns joined name parts",
			parts: []string{" Johann", "     ", "   ", "	Sebastian  ", "Bach"},
			want: "Johann Sebastian Bach",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := name.FullName(fallback, tC.parts...)
			if got != tC.want {
				t.Errorf("FullName(%s, %v) = %s, want %s", fallback, tC.parts, got, tC.want)
			}
		})
	}
}
