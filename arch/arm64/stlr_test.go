package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStlrCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"stlr x0,[x1]",
			ctorStlr(t, xreg(t, 0), xreg(t, 1)),
			0xc89ffc20,
		},
		{
			"stlr w2,[sp]",
			ctorStlr(t, wreg(t, 2), SP),
			0x889fffe2,
		},
		{
			"stlr xzr,[x3]",
			ctorStlr(t, XZR, xreg(t, 3)),
			0xc89ffc7f,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorStlr(t, xreg(t, 0), xreg(t, 1))
	_, ok := in.(Stlr)
	require.True(t, ok, "type = %T, want Stlr", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"stlr sp,rt",
			func() error {
				_, err := New().Stlr(SP, xreg(t, 1))
				return err
			},
		},
		{
			"stlr w1 base",
			func() error {
				_, err := New().Stlr(xreg(t, 0), wreg(t, 1))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorStlr — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorStlr(t *testing.T, rt, rn Reg) Instr {
	t.Helper()
	in, err := New().Stlr(rt, rn)
	require.NoError(t, err)
	return in
}
