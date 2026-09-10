package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddImmBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		imm  int64
		sh   Sh12
		word uint32
	}{
		{"add x0,x1,#0x42", "x0", "x1", 0x42, NoSh12, 0x91010820},
		{"add x0,x1,#1,lsl#12", "x0", "x1", 1, LSL12, 0x91400420},
		{"add w2,w3,#7", "w2", "w3", 7, NoSh12, 0x11001c62},
		{"add sp,sp,#0x10", "sp", "sp", 0x10, NoSh12, 0x910043ff},
	}
	for _, c := range cases {
		in, err := New().AddImm(reg(t, c.rd), reg(t, c.rn), imm12(t, c.imm), c.sh)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().AddImm(reg(t, first.rd), reg(t, first.rn), imm12(t, first.imm), first.sh)
	require.NoError(t, err)
	_, ok := in.(AddImm)
	require.True(t, ok, "type = %T, want AddImm", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		imm  int64
		sh   Sh12
	}{
		{"add xzr", "xzr", "x1", 1, NoSh12},
		{"add rn xzr", "x0", "xzr", 1, NoSh12},
		{"add sp+w", "sp", "w1", 1, NoSh12},
	}
	for _, c := range errCases {
		_, err := New().AddImm(reg(t, c.rd), reg(t, c.rn), imm12(t, c.imm), c.sh)
		assertErr(t, c.name, err)
	}
}
