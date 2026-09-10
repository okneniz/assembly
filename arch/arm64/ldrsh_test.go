package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdrshBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"ldrsh x0,[x1]", "x0", "x1", 0, 0x79800020},
		{"ldrsh xzr,[x2,#0x1ffe]", "xzr", "x2", 0x1ffe, 0x79bffc5f},
	}
	for _, c := range cases {
		in, err := New().Ldrsh(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldrsh(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Ldrsh)
	require.True(t, ok, "type = %T, want Ldrsh", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
		off  Off
	}{
		{"ldrsh w,rt", "w0", "x1", 0},
		{"ldrsh w1 base", "x0", "w1", 0},
		{"ldrsh off 1", "x0", "x1", 1},
	}
	for _, c := range errCases {
		_, err := New().Ldrsh(reg(t, c.rt), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
