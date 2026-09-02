package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRevCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"rev x0,x1",
			ctorRev(t, xreg(t, 0), xreg(t, 1)),
			0xdac00c20,
		},
		{
			"rev xzr,x2",
			ctorRev(t, XZR, xreg(t, 2)),
			0xdac00c5f,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorRev(t, xreg(t, 0), xreg(t, 1))
	_, ok := in.(Rev)
	require.True(t, ok, "type = %T, want Rev", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"rev w form",
			func() error {
				_, err := New().Rev(wreg(t, 0), wreg(t, 1))
				return err
			},
		},
		{
			"rev w,rn",
			func() error {
				_, err := New().Rev(xreg(t, 0), wreg(t, 1))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorRev — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorRev(t *testing.T, rd, rn Reg) Instr {
	t.Helper()
	in, err := New().Rev(rd, rn)
	require.NoError(t, err)
	return in
}
