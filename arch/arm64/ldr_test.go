package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdrBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"ldr x0,[x1]", "x0", "x1", 0, 0xf9400020},
		{"ldr x0,[x1,#8]", "x0", "x1", 8, 0xf9400420},
		{"ldr w2,[sp,#0x10]", "w2", "sp", 0x10, 0xb94013e2},
		{"ldr xzr,[x1,#0x7ff8]", "xzr", "x1", 0x7ff8, 0xf97ffc3f},
	}
	for _, c := range cases {
		in, err := New().Ldr(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldr(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Ldr)
	require.True(t, ok, "type = %T, want Ldr", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
		off  Off
	}{
		{"ldr sp,rt", "sp", "x1", 0},
		{"ldr w,[w1]", "w0", "w1", 0},
		{"ldr xzr base", "x0", "xzr", 0},
		{"ldr off 4 in x form", "x0", "x1", 4},
		{"ldr off 0x8000 in w form", "w0", "x1", 0x8000},
	}
	for _, c := range errCases {
		_, err := New().Ldr(reg(t, c.rt), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
