package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdarBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		word uint32
	}{
		{"ldar x0,[x1]", "x0", "x1", 0xc8dffc20},
		{"ldar w2,[sp]", "w2", "sp", 0x88dfffe2},
		{"ldar xzr,[x3]", "xzr", "x3", 0xc8dffc7f},
	}
	for _, c := range cases {
		in, err := New().Ldar(reg(t, c.rt), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldar(reg(t, first.rt), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Ldar)
	require.True(t, ok, "type = %T, want Ldar", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
	}{
		{"ldar sp,rt", "sp", "x1"},
		{"ldar w1 base", "x0", "w1"},
	}
	for _, c := range errCases {
		_, err := New().Ldar(reg(t, c.rt), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
