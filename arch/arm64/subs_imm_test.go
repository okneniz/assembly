package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubsImmBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		imm  int64
		sh   Sh12
		word uint32
	}{
		{"subs x0,x1,#0x42", "x0", "x1", 0x42, NoSh12, 0xf1010820},
		{"subs xzr,x1,#1 (cmp)", "xzr", "x1", 1, NoSh12, 0xf100043f},
		{"subs w2,w3,#7,lsl#12", "w2", "w3", 7, LSL12, 0x71401c62},
		{"subs x2,sp,#0x10", "x2", "sp", 0x10, NoSh12, 0xf10043e2},
	}
	for _, c := range cases {
		in, err := New().SubsImm(reg(t, c.rd), reg(t, c.rn), imm12(t, c.imm), c.sh)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().SubsImm(reg(t, first.rd), reg(t, first.rn), imm12(t, first.imm), first.sh)
	require.NoError(t, err)
	_, ok := in.(SubsImm)
	require.True(t, ok, "type = %T, want SubsImm", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		imm  int64
		sh   Sh12
	}{
		{"subs sp,rd", "sp", "x1", 1, NoSh12},
		{"subs xzr,rn", "x0", "xzr", 1, NoSh12},
		{"subs x,w widths", "x0", "w1", 1, NoSh12},
	}
	for _, c := range errCases {
		_, err := New().SubsImm(reg(t, c.rd), reg(t, c.rn), imm12(t, c.imm), c.sh)
		assertErr(t, c.name, err)
	}
}
