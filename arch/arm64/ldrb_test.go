package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdrbBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"ldrb w0,[x1]", "w0", "x1", 0, 0x39400020},
		{"ldrb w2,[sp,#1]", "w2", "sp", 1, 0x394007e2},
		{"ldrb wzr,[x2,#0xfff]", "wzr", "x2", 0xfff, 0x397ffc5f},
	}
	for _, c := range cases {
		in, err := New().Ldrb(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldrb(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Ldrb)
	require.True(t, ok, "type = %T, want Ldrb", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
		off  Off
	}{
		{"ldrb x,rt", "x0", "x1", 0},
		{"ldrb w1 base", "w0", "w1", 0},
		{"ldrb off -1", "w0", "x1", -1},
		{"ldrb off 0x1000", "w0", "x1", 0x1000},
	}
	for _, c := range errCases {
		_, err := New().Ldrb(reg(t, c.rt), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
