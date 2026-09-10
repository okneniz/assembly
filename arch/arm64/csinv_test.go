package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCsinvBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		cond string
		word uint32
	}{
		{"csinv x5,x6,x7,al", "x5", "x6", "x7", "al", 0xda87e0c5},
		{"csinv w2,w3,wzr,le", "w2", "w3", "wzr", "le", 0x5a9fd062},
	}
	for _, c := range cases {
		in, err := New().Csinv(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.cond)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Csinv(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), first.cond)
	require.NoError(t, err)
	_, ok := in.(Csinv)
	require.True(t, ok, "type = %T, want Csinv", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		cond string
	}{
		{"csinv sp,rd", "sp", "x6", "x7", "al"},
		{"csinv w,x widths", "x5", "w6", "x7", "al"},
		{"csinv bad cond", "x5", "x6", "x7", "foo"},
	}
	for _, c := range errCases {
		_, err := New().Csinv(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.cond)
		assertErr(t, c.name, err)
	}
}
