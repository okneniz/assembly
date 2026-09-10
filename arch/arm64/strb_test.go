package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStrbBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"strb w0,[x1]", "w0", "x1", 0, 0x39000020},
		{"strb wzr,[sp,#1]", "wzr", "sp", 1, 0x390007ff},
		{"strb w2,[x3,#0xfff]", "w2", "x3", 0xfff, 0x393ffc62},
	}
	for _, c := range cases {
		in, err := New().Strb(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Strb(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Strb)
	require.True(t, ok, "type = %T, want Strb", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
		off  Off
	}{
		{"strb x,rt", "x0", "x1", 0},
		{"strb w1 base", "w0", "w1", 0},
		{"strb off -1", "w0", "x1", -1},
	}
	for _, c := range errCases {
		_, err := New().Strb(reg(t, c.rt), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
