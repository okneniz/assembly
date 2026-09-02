package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestTbzCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"tbz w0,#5,0x1000",
			ctorTbz(t, wreg(t, 0), 5, 0x1000),
			0x36280000,
		},
		{
			"tbz w1,#0,0x1010",
			ctorTbz(t, wreg(t, 1), 0, 0x1010),
			0x36000081,
		},
		{
			"tbz w0,#31,0x8ffc",
			ctorTbz(t, wreg(t, 0), 31, 0x8ffc),
			0x36fbffe0,
		},
		{
			"tbz x2,#32,0x1000",
			ctorTbz(t, xreg(t, 2), 32, 0x1000),
			0xb6000002,
		},
		{
			"tbz w0,#1,0xf00",
			ctorTbz(t, wreg(t, 0), 1, 0xf00),
			0x360ff800,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorTbz(t, wreg(t, 0), 5, 0x1000)
	_, ok := in.(Tbz)
	require.True(t, ok, "type = %T, want Tbz", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"tbz sp,rt",
			func() error {
				_, err := New().Tbz(SP, 5, 0x1000)
				return err
			},
		},
		{
			"tbz bit 64",
			func() error {
				_, err := New().Tbz(xreg(t, 0), 64, 0x1000)
				return err
			},
		},
		{
			"tbz x0,#5",
			func() error {
				_, err := New().Tbz(xreg(t, 0), 5, 0x1000)
				return err
			},
		},
		{
			"tbz w0,#32",
			func() error {
				_, err := New().Tbz(wreg(t, 0), 32, 0x1000)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorTbz — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorTbz(t *testing.T, rt Reg, bit uint32, target int64) Instr {
	t.Helper()
	in, err := New().Tbz(rt, bit, target)
	require.NoError(t, err)
	return in
}
