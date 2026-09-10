package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSturbBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"sturb w0,[x1]", "w0", "x1", 0, 0x38000020},
		{"sturb wzr,[x29,#255]", "wzr", "x29", 255, 0x380ff3bf},
	}
	for _, c := range cases {
		in, err := New().Sturb(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Sturb(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Sturb)
	require.True(t, ok, "type = %T, want Sturb", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
		off  Off
	}{
		{"sturb x,rt", "x0", "x1", 0},
		{"sturb w1 base", "w0", "w1", 0},
		{"sturb off 256", "w0", "x1", 256},
	}
	for _, c := range errCases {
		_, err := New().Sturb(reg(t, c.rt), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
