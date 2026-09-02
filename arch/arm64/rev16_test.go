package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRev16Ctor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"rev16 x0,x1",
			ctorRev16(t, xreg(t, 0), xreg(t, 1)),
			0xdac00420,
		},
		{
			"rev16 w2,w3",
			ctorRev16(t, wreg(t, 2), wreg(t, 3)),
			0x5ac00462,
		},
		{
			"rev16 x1,xzr",
			ctorRev16(t, xreg(t, 1), XZR),
			0xdac007e1,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorRev16(t, xreg(t, 0), xreg(t, 1))
	_, ok := in.(Rev16)
	require.True(t, ok, "type = %T, want Rev16", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"rev16 sp,rd",
			func() error {
				_, err := New().Rev16(SP, xreg(t, 1))
				return err
			},
		},
		{
			"rev16 sp,rn",
			func() error {
				_, err := New().Rev16(xreg(t, 0), SP)
				return err
			},
		},
		{
			"rev16 x,w widths",
			func() error {
				_, err := New().Rev16(xreg(t, 0), wreg(t, 1))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorRev16 — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorRev16(t *testing.T, rd, rn Reg) Instr {
	t.Helper()
	in, err := New().Rev16(rd, rn)
	require.NoError(t, err)
	return in
}
