package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestMovnBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"movn x0,#0",
			buildMovn(t, xreg(t, 0), imm16(t, 0), Hw0),
			0x92800000,
		},
		{
			"movn x1,#0x1234,Hw1",
			buildMovn(t, xreg(t, 1), imm16(t, 0x1234), Hw1),
			0x92a24681,
		},
		{
			"movn w2,#0xffff",
			buildMovn(t, wreg(t, 2), imm16(t, 0xffff), Hw0),
			0x129fffe2,
		},
		{
			"movn x0,#0xffff,Hw3",
			buildMovn(t, xreg(t, 0), imm16(t, 0xffff), Hw3),
			0x92ffffe0,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := buildMovn(t, xreg(t, 0), imm16(t, 0), Hw0)
	_, ok := in.(Movn)
	require.True(t, ok, "type = %T, want Movn", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"movn sp,rd",
			func() error {
				_, err := New().Movn(SP, imm16(t, 0), Hw0)
				return err
			},
		},
		{
			"movn w,Hw2",
			func() error {
				_, err := New().Movn(wreg(t, 0), imm16(t, 0), Hw2)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildMovn — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildMovn(t *testing.T, rd Reg, imm Imm16, hw Hw) Instr {
	t.Helper()
	in, err := New().Movn(rd, imm, hw)
	require.NoError(t, err)
	return in
}
