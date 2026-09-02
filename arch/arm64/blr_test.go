package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBlrCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"blr x1",
			ctorBlr(t, xreg(t, 1)),
			0xd63f0020,
		},
		{
			"blr x30",
			ctorBlr(t, xreg(t, 30)),
			0xd63f03c0,
		},
		{
			"blr xzr",
			ctorBlr(t, XZR),
			0xd63f03e0,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorBlr(t, xreg(t, 1))
	_, ok := in.(Blr)
	require.True(t, ok, "type = %T, want Blr", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"blr w1",
			func() error {
				_, err := New().Blr(wreg(t, 1))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorBlr — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorBlr(t *testing.T, rn Reg) Instr {
	t.Helper()
	in, err := New().Blr(rn)
	require.NoError(t, err)
	return in
}
