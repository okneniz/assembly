package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStlrbCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"stlrb w3,[x4]",
			ctorStlrb(t, wreg(t, 3), xreg(t, 4)),
			0x089ffc83,
		},
		{
			"stlrb wzr,[sp]",
			ctorStlrb(t, WZR, SP),
			0x089fffff,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorStlrb(t, wreg(t, 3), xreg(t, 4))
	_, ok := in.(Stlrb)
	require.True(t, ok, "type = %T, want Stlrb", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"stlrb x,rt",
			func() error {
				_, err := New().Stlrb(xreg(t, 3), xreg(t, 4))
				return err
			},
		},
		{
			"stlrb w1 base",
			func() error {
				_, err := New().Stlrb(wreg(t, 3), wreg(t, 4))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorStlrb — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorStlrb(t *testing.T, rt, rn Reg) Instr {
	t.Helper()
	in, err := New().Stlrb(rt, rn)
	require.NoError(t, err)
	return in
}
