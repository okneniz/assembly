package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdaxrBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		word uint32
	}{
		{"ldaxr x0,[x1]", "x0", "x1", 0xc85ffc20},
		{"ldaxr w1,[x2]", "w1", "x2", 0x885ffc41},
		{"ldaxr w2,[sp]", "w2", "sp", 0x885fffe2},
	}
	for _, c := range cases {
		in, err := New().Ldaxr(reg(t, c.rt), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldaxr(reg(t, first.rt), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Ldaxr)
	require.True(t, ok, "type = %T, want Ldaxr", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
	}{
		{"ldaxr sp,rt", "sp", "x1"},
		{"ldaxr w1 base", "x0", "w1"},
	}
	for _, c := range errCases {
		_, err := New().Ldaxr(reg(t, c.rt), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
