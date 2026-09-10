package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubsShiftBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
		word uint32
	}{
		{"subs x0,x1,x2", "x0", "x1", "x2", 0, LSL, 0xeb020020},
		{"subs x0,x1,x2,lsl#3", "x0", "x1", "x2", 3, LSL, 0xeb020c20},
		{"subs xzr,x1,x2,asr#4 (cmp)", "xzr", "x1", "x2", 4, ASR, 0xeb82103f},
		{"subs w3,w4,w5,lsl#2", "w3", "w4", "w5", 2, LSL, 0x6b050883},
	}
	for _, c := range cases {
		in, err := New().SubsShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().SubsShift(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), imm6(t, first.imm), first.sh)
	require.NoError(t, err)
	_, ok := in.(SubsShift)
	require.True(t, ok, "type = %T, want SubsShift", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		imm  int64
		sh   Shift
	}{
		{"subs sp,rd", "sp", "x1", "x2", 0, LSL},
		{"subs sp,rn", "x0", "sp", "x2", 0, LSL},
		{"subs sp,rm", "x0", "x1", "sp", 0, LSL},
		{"subs x,w widths", "x0", "w1", "x2", 0, LSL},
		{"subs ror shift", "x0", "x1", "x2", 0, ROR},
		{"subs w,#32 shift", "w0", "w1", "w2", 32, LSL},
	}
	for _, c := range errCases {
		_, err := New().SubsShift(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.imm), c.sh)
		assertErr(t, c.name, err)
	}
}
