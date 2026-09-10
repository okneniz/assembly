package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdaxrbBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		word uint32
	}{
		{"ldaxrb w0,[x1]", "w0", "x1", 0x085ffc20},
		{"ldaxrb w5,[sp]", "w5", "sp", 0x085fffe5},
	}
	for _, c := range cases {
		in, err := New().Ldaxrb(reg(t, c.rt), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldaxrb(reg(t, first.rt), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Ldaxrb)
	require.True(t, ok, "type = %T, want Ldaxrb", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
	}{
		{"ldaxrb x,rt", "x0", "x1"},
		{"ldaxrb w1 base", "w0", "w1"},
	}
	for _, c := range errCases {
		_, err := New().Ldaxrb(reg(t, c.rt), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
