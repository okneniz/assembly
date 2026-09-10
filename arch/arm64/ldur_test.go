package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLdurBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"ldur x0,[x1]", "x0", "x1", 0, 0xf8400020},
		{"ldur x1,[x2,#-1]", "x1", "x2", -1, 0xf85ff041},
		{"ldur w2,[sp,#255]", "w2", "sp", 255, 0xb84ff3e2},
		{"ldur w3,[x4,#-256]", "w3", "x4", -256, 0xb8500083},
	}
	for _, c := range cases {
		in, err := New().Ldur(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Ldur(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Ldur)
	require.True(t, ok, "type = %T, want Ldur", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
		off  Off
	}{
		{"ldur sp,rt", "sp", "x1", 0},
		{"ldur w1 base", "x0", "w1", 0},
		{"ldur off 256", "x0", "x1", 256},
	}
	for _, c := range errCases {
		_, err := New().Ldur(reg(t, c.rt), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
