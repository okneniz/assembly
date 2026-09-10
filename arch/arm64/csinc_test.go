package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCsincBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		cond string
		word uint32
	}{
		{"csinc x1,x2,x3,eq", "x1", "x2", "x3", "eq", 0x9a830441},
		{"csinc w0,wzr,w1,ne", "w0", "wzr", "w1", "ne", 0x1a8117e0},
	}
	for _, c := range cases {
		in, err := New().Csinc(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.cond)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Csinc(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), first.cond)
	require.NoError(t, err)
	_, ok := in.(Csinc)
	require.True(t, ok, "type = %T, want Csinc", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		cond string
	}{
		{"csinc sp,rd", "sp", "x2", "x3", "eq"},
		{"csinc w,x widths", "x1", "w2", "x3", "eq"},
		{"csinc bad cond", "x1", "x2", "x3", "foo"},
	}
	for _, c := range errCases {
		_, err := New().Csinc(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.cond)
		assertErr(t, c.name, err)
	}
}
