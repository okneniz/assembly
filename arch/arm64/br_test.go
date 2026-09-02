package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBrCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"br x5",
			ctorBr(t, xreg(t, 5)),
			0xd61f00a0,
		},
		{
			"br x0",
			ctorBr(t, xreg(t, 0)),
			0xd61f0000,
		},
		{
			"br xzr",
			ctorBr(t, XZR),
			0xd61f03e0,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorBr(t, xreg(t, 5))
	_, ok := in.(Br)
	require.True(t, ok, "type = %T, want Br", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"br w5",
			func() error {
				_, err := New().Br(wreg(t, 5))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorBr — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorBr(t *testing.T, rn Reg) Instr {
	t.Helper()
	in, err := New().Br(rn)
	require.NoError(t, err)
	return in
}
