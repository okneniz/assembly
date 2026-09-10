package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCcmpBuild(t *testing.T) {
	cases := []struct {
		name string
		rn   string
		rm   string
		nzcv uint32
		cond string
		word uint32
	}{
		{"ccmp x1,x2,#3,eq", "x1", "x2", 3, "eq", 0xfa420023},
		{"ccmp x0,xzr,#0xf,ne", "x0", "xzr", 0xf, "ne", 0xfa5f100f},
	}
	for _, c := range cases {
		in, err := New().Ccmp(reg(t, c.rn), reg(t, c.rm), c.nzcv, c.cond)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ccmp(reg(t, first.rn), reg(t, first.rm), first.nzcv, first.cond)
	require.NoError(t, err)
	_, ok := in.(Ccmp)
	require.True(t, ok, "type = %T, want Ccmp", in)

	errCases := []struct {
		name string
		rn   string
		rm   string
		nzcv uint32
		cond string
	}{
		{"ccmp w form", "w1", "x2", 3, "eq"},
		{"ccmp rm w form", "x1", "w2", 3, "eq"},
		{"ccmp nzcv=0x10", "x1", "x2", 0x10, "eq"},
		{"ccmp bad cond", "x1", "x2", 3, "foo"},
	}
	for _, c := range errCases {
		_, err := New().Ccmp(reg(t, c.rn), reg(t, c.rm), c.nzcv, c.cond)
		assertErr(t, c.name, err)
	}
}
