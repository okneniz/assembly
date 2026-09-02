package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRev32Ctor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"rev32 x0,x1",
			ctorRev32(t, xreg(t, 0), xreg(t, 1)),
			0xdac00820,
		},
		{
			"rev32 x4,xzr",
			ctorRev32(t, xreg(t, 4), XZR),
			0xdac00be4,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorRev32(t, xreg(t, 0), xreg(t, 1))
	_, ok := in.(Rev32)
	require.True(t, ok, "type = %T, want Rev32", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"rev32 w form",
			func() error {
				_, err := New().Rev32(wreg(t, 0), wreg(t, 1))
				return err
			},
		},
		{
			"rev32 w,rn",
			func() error {
				_, err := New().Rev32(xreg(t, 0), wreg(t, 1))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorRev32 — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorRev32(t *testing.T, rd, rn Reg) Instr {
	t.Helper()
	in, err := New().Rev32(rd, rn)
	require.NoError(t, err)
	return in
}
