package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCsnegBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		cond string
		word uint32
	}{
		{"csneg x0,x1,x2,mi", "x0", "x1", "x2", "mi", 0xda824420},
		{"csneg w4,wzr,wzr,eq", "w4", "wzr", "wzr", "eq", 0x5a9f07e4},
	}
	for _, c := range cases {
		in, err := New().Csneg(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.cond)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Csneg(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), first.cond)
	require.NoError(t, err)
	_, ok := in.(Csneg)
	require.True(t, ok, "type = %T, want Csneg", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		cond string
	}{
		{"csneg sp,rd", "sp", "x1", "x2", "mi"},
		{"csneg w,x widths", "x0", "w1", "x2", "mi"},
		{"csneg bad cond", "x0", "x1", "x2", "foo"},
	}
	for _, c := range errCases {
		_, err := New().Csneg(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), c.cond)
		assertErr(t, c.name, err)
	}
}
