package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdaxrCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldaxr x0,[x1]",
			ctorLdaxr(t, xreg(t, 0), xreg(t, 1)),
			0xc85ffc20,
		},
		{
			"ldaxr w1,[x2]",
			ctorLdaxr(t, wreg(t, 1), xreg(t, 2)),
			0x885ffc41,
		},
		{
			"ldaxr w2,[sp]",
			ctorLdaxr(t, wreg(t, 2), SP),
			0x885fffe2,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorLdaxr(t, xreg(t, 0), xreg(t, 1))
	_, ok := in.(Ldaxr)
	require.True(t, ok, "type = %T, want Ldaxr", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldaxr sp,rt",
			func() error {
				_, err := New().Ldaxr(SP, xreg(t, 1))
				return err
			},
		},
		{
			"ldaxr w1 base",
			func() error {
				_, err := New().Ldaxr(xreg(t, 0), wreg(t, 1))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorLdaxr — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorLdaxr(t *testing.T, rt, rn Reg) Instr {
	t.Helper()
	in, err := New().Ldaxr(rt, rn)
	require.NoError(t, err)
	return in
}
