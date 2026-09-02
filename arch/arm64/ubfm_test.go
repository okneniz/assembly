package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestUbfmCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ubfm x0,x1,#5,#63 (lsr)",
			ctorUbfm(t, xreg(t, 0), xreg(t, 1), 5, 63),
			0xd345fc20,
		},
		{
			"ubfm x2,x3,#59,#58 (lsl)",
			ctorUbfm(t, xreg(t, 2), xreg(t, 3), 59, 58),
			0xd37be862,
		},
		{
			"ubfm w1,w2,#0,#31 (lsr w)",
			ctorUbfm(t, wreg(t, 1), wreg(t, 2), 0, 31),
			0x53007c41,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorUbfm(t, xreg(t, 0), xreg(t, 1), 5, 63)
	_, ok := in.(Ubfm)
	require.True(t, ok, "type = %T, want Ubfm", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ubfm sp,rd",
			func() error {
				_, err := New().Ubfm(SP, xreg(t, 1), 5, 63)
				return err
			},
		},
		{
			"ubfm sp,rn",
			func() error {
				_, err := New().Ubfm(xreg(t, 0), SP, 5, 63)
				return err
			},
		},
		{
			"ubfm x,w widths",
			func() error {
				_, err := New().Ubfm(xreg(t, 0), wreg(t, 1), 5, 63)
				return err
			},
		},
		{
			"ubfm immr 64",
			func() error {
				_, err := New().Ubfm(xreg(t, 0), xreg(t, 1), 64, 63)
				return err
			},
		},
		{
			"ubfm imms 64",
			func() error {
				_, err := New().Ubfm(xreg(t, 0), xreg(t, 1), 5, 64)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorUbfm — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorUbfm(t *testing.T, rd, rn Reg, immr, imms uint32) Instr {
	t.Helper()
	in, err := New().Ubfm(rd, rn, immr, imms)
	require.NoError(t, err)
	return in
}
