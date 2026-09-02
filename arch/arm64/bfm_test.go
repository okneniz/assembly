package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBfmCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"bfm x0,x1,#5,#7",
			ctorBfm(t, xreg(t, 0), xreg(t, 1), 5, 7),
			0xb3451c20,
		},
		{
			"bfm w2,w3,#0,#31",
			ctorBfm(t, wreg(t, 2), wreg(t, 3), 0, 31),
			0x33007c62,
		},
		{
			"bfm x0,x1,#63,#63",
			ctorBfm(t, xreg(t, 0), xreg(t, 1), 63, 63),
			0xb37ffc20,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorBfm(t, xreg(t, 0), xreg(t, 1), 5, 7)
	_, ok := in.(Bfm)
	require.True(t, ok, "type = %T, want Bfm", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"bfm sp,rd",
			func() error {
				_, err := New().Bfm(SP, xreg(t, 1), 5, 7)
				return err
			},
		},
		{
			"bfm sp,rn",
			func() error {
				_, err := New().Bfm(xreg(t, 0), SP, 5, 7)
				return err
			},
		},
		{
			"bfm x,w widths",
			func() error {
				_, err := New().Bfm(xreg(t, 0), wreg(t, 1), 5, 7)
				return err
			},
		},
		{
			"bfm immr 64",
			func() error {
				_, err := New().Bfm(xreg(t, 0), xreg(t, 1), 64, 7)
				return err
			},
		},
		{
			"bfm imms 64",
			func() error {
				_, err := New().Bfm(xreg(t, 0), xreg(t, 1), 5, 64)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorBfm — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorBfm(t *testing.T, rd, rn Reg, immr, imms uint32) Instr {
	t.Helper()
	in, err := New().Bfm(rd, rn, immr, imms)
	require.NoError(t, err)
	return in
}
