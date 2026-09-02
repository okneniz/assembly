package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdarbCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldarb w3,[x4]",
			ctorLdarb(t, wreg(t, 3), xreg(t, 4)),
			0x08dffc83,
		},
		{
			"ldarb wzr,[x29]",
			ctorLdarb(t, WZR, xreg(t, 29)),
			0x08dfffbf,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorLdarb(t, wreg(t, 3), xreg(t, 4))
	_, ok := in.(Ldarb)
	require.True(t, ok, "type = %T, want Ldarb", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldarb x,rt",
			func() error {
				_, err := New().Ldarb(xreg(t, 3), xreg(t, 4))
				return err
			},
		},
		{
			"ldarb w1 base",
			func() error {
				_, err := New().Ldarb(wreg(t, 3), wreg(t, 4))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorLdarb — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorLdarb(t *testing.T, rt, rn Reg) Instr {
	t.Helper()
	in, err := New().Ldarb(rt, rn)
	require.NoError(t, err)
	return in
}
