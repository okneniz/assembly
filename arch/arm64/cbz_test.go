package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCbzCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"cbz x0,0x1000",
			ctorCbz(t, xreg(t, 0), 0x1000),
			0xb4000000,
		},
		{
			"cbz w3,0x1014",
			ctorCbz(t, wreg(t, 3), 0x1014),
			0x340000a3,
		},
		{
			"cbz x1,0xfe8",
			ctorCbz(t, xreg(t, 1), 0xfe8),
			0xb4ffff41,
		},
		{
			"cbz x0,0x800",
			ctorCbz(t, xreg(t, 0), 0x800),
			0xb4ffc000,
		},
		{
			"cbz x0,0x100ffc",
			ctorCbz(t, xreg(t, 0), 0x100ffc),
			0xb47fffe0,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorCbz(t, xreg(t, 0), 0x1000)
	_, ok := in.(Cbz)
	require.True(t, ok, "type = %T, want Cbz", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"cbz sp,rt",
			func() error {
				_, err := New().Cbz(SP, 0x1000)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorCbz — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorCbz(t *testing.T, rt Reg, target int64) Instr {
	t.Helper()
	in, err := New().Cbz(rt, target)
	require.NoError(t, err)
	return in
}
