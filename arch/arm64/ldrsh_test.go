package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdrshCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldrsh x0,[x1]",
			ctorLdrsh(t, xreg(t, 0), xreg(t, 1), 0),
			0x79800020,
		},
		{
			"ldrsh xzr,[x2,#0x1ffe]",
			ctorLdrsh(t, XZR, xreg(t, 2), 0x1ffe),
			0x79bffc5f,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorLdrsh(t, xreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Ldrsh)
	require.True(t, ok, "type = %T, want Ldrsh", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldrsh w,rt",
			func() error {
				_, err := New().Ldrsh(wreg(t, 0), xreg(t, 1), 0)
				return err
			},
		},
		{
			"ldrsh w1 base",
			func() error {
				_, err := New().Ldrsh(xreg(t, 0), wreg(t, 1), 0)
				return err
			},
		},
		{
			"ldrsh off 1",
			func() error {
				_, err := New().Ldrsh(xreg(t, 0), xreg(t, 1), 1)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorLdrsh — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorLdrsh(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Ldrsh(rt, rn, off)
	require.NoError(t, err)
	return in
}
