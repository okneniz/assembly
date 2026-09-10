package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStlrbBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		word uint32
	}{
		{"stlrb w3,[x4]", "w3", "x4", 0x089ffc83},
		{"stlrb wzr,[sp]", "wzr", "sp", 0x089fffff},
	}
	for _, c := range cases {
		in, err := New().Stlrb(reg(t, c.rt), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Stlrb(reg(t, first.rt), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Stlrb)
	require.True(t, ok, "type = %T, want Stlrb", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
	}{
		{"stlrb x,rt", "x3", "x4"},
		{"stlrb w1 base", "w3", "w4"},
	}
	for _, c := range errCases {
		_, err := New().Stlrb(reg(t, c.rt), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
