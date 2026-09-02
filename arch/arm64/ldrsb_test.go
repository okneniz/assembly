package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdrsbCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldrsb x0,[x1]",
			ctorLdrsb(t, xreg(t, 0), xreg(t, 1), 0),
			0x39800020,
		},
		{
			"ldrsb x2,[sp,#0xfff]",
			ctorLdrsb(t, xreg(t, 2), SP, 0xfff),
			0x39bfffe2,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorLdrsb(t, xreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Ldrsb)
	require.True(t, ok, "type = %T, want Ldrsb", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldrsb w,rt",
			func() error {
				_, err := New().Ldrsb(wreg(t, 0), xreg(t, 1), 0)
				return err
			},
		},
		{
			"ldrsb w1 base",
			func() error {
				_, err := New().Ldrsb(xreg(t, 0), wreg(t, 1), 0)
				return err
			},
		},
		{
			"ldrsb off -1",
			func() error {
				_, err := New().Ldrsb(xreg(t, 0), xreg(t, 1), -1)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorLdrsb — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorLdrsb(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Ldrsb(rt, rn, off)
	require.NoError(t, err)
	return in
}
