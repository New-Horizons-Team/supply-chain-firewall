package slice_utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test ContainsPrefix.
func TestContainsPrefix(t *testing.T) {
	type testCase struct {
		s    []string
		v    string
		sep  string
		want string
	}

	testCases := []testCase{{
		s:    []string{"foo", "baz"},
		v:    "foo",
		sep:  "|",
		want: "foo",
	}, {
		s:    []string{"foo", "baz"},
		v:    "foobar",
		sep:  "|",
		want: "",
	}, {
		s:    []string{"foo", "baz"},
		v:    "foo|bar",
		sep:  "|",
		want: "foo",
	}, {
		s:    []string{"foo", "baz"},
		v:    "baz",
		sep:  "|",
		want: "baz",
	}, {
		s:    []string{"-o", "-orb"},
		v:    "-o",
		sep:  "=",
		want: "-o",
	}, {
		s:    []string{"-o", "-orb"},
		v:    "-ops",
		sep:  "=",
		want: "",
	}, {
		s:    []string{"-o", "-orb"},
		v:    "-o=foo",
		sep:  "=",
		want: "-o",
	}, {
		s:    []string{"-o", "-orb"},
		v:    "-orb",
		sep:  "=",
		want: "-orb",
	}, {
		s:    []string{"-o", "-orb"},
		v:    "-orb=asd",
		sep:  "=",
		want: "-orb",
	}}

	for _, tc := range testCases {
		got, ok := ContainsPrefix(tc.s, tc.v, tc.sep)
		require.Equal(t, tc.want != "", ok)
		require.Equal(t, tc.want, got)
	}
}
