package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdrswCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldrsw x0,[x1]",
			ctorLdrsw(t, xreg(t, 0), xreg(t, 1), 0),
			0xb9800020,
		},
		{
			"ldrsw x2,[sp,#0xffc]",
			ctorLdrsw(t, xreg(t, 2), SP, 0xffc),
			0xb98fffe2,
		},
		{
			"ldrsw xzr,[x3,#4]",
			ctorLdrsw(t, XZR, xreg(t, 3), 4),
			0xb980047f,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorLdrsw(t, xreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Ldrsw)
	require.True(t, ok, "type = %T, want Ldrsw", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldrsw w,rt",
			func() error {
				_, err := New().Ldrsw(wreg(t, 0), xreg(t, 1), 0)
				return err
			},
		},
		{
			"ldrsw w1 base",
			func() error {
				_, err := New().Ldrsw(xreg(t, 0), wreg(t, 1), 0)
				return err
			},
		},
		{
			"ldrsw off 2",
			func() error {
				_, err := New().Ldrsw(xreg(t, 0), xreg(t, 1), 2)
				return err
			},
		},
		{
			"ldrsw off 0x10000",
			func() error {
				_, err := New().Ldrsw(xreg(t, 0), xreg(t, 1), 0x10000)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorLdrsw — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorLdrsw(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Ldrsw(rt, rn, off)
	require.NoError(t, err)
	return in
}
