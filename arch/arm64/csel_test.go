package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCselBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		cond string
		word uint32
	}{
		{"csel x0,x1,x2,eq", "x0", "x1", "x2", "eq", 0x9a820020},
		{"csel w3,w4,w5,hs", "w3", "w4", "w5", "hs", 0x1a852083},
		{"csel x0,xzr,x1,gt", "x0", "xzr", "x1", "gt", 0x9a81c3e0},
	}
	for _, c := range cases {
		in, err := New().Csel(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.cond)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Csel(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), first.cond)
	require.NoError(t, err)
	_, ok := in.(Csel)
	require.True(t, ok, "type = %T, want Csel", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		cond string
	}{
		{"csel sp,rd", "sp", "x1", "x2", "eq"},
		{"csel sp,rn", "x0", "sp", "x2", "eq"},
		{"csel sp,rm", "x0", "x1", "sp", "eq"},
		{"csel w,x widths", "x0", "w1", "x2", "eq"},
		{"csel bad cond", "x0", "x1", "x2", "foo"},
	}
	for _, c := range errCases {
		_, err := New().Csel(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.cond)
		assertErr(t, c.name, err)
	}
}
