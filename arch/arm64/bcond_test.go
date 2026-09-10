package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBcondBuild(t *testing.T) {
	cases := []struct {
		name   string
		cond   string
		target int64
		word   uint32
	}{
		{"b.eq 0x1000", "eq", 0, 0x54000000},
		{"b.ne 0x1010", "ne", 16, 0x54000081},
		{"b.hs 0xffc", "hs", -4, 0x54ffffe2},
		{"b.eq 0x100ffc", "eq", 1048572, 0x547fffe0},
	}
	for _, c := range cases {
		in, err := New().Bcond(c.cond, c.target)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Bcond(first.cond, first.target)
	require.NoError(t, err)
	_, ok := in.(Bcond)
	require.True(t, ok, "type = %T, want Bcond", in)

	errCases := []struct {
		name   string
		cond   string
		target int64
	}{
		{"b.cond bad cond", "foo", 0x1000},
	}
	for _, c := range errCases {
		_, err := New().Bcond(c.cond, c.target)
		assertErr(t, c.name, err)
	}
}
