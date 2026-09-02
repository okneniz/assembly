package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSbfmCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"sbfm x0,x1,#5,#63 (asr)",
			ctorSbfm(t, xreg(t, 0), xreg(t, 1), 5, 63),
			0x9345fc20,
		},
		{
			"sbfm x1,x2,#0,#7 (sxtb)",
			ctorSbfm(t, xreg(t, 1), xreg(t, 2), 0, 7),
			0x93401c41,
		},
		{
			"sbfm w3,w4,#1,#3 (sbfx)",
			ctorSbfm(t, wreg(t, 3), wreg(t, 4), 1, 3),
			0x13010c83,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorSbfm(t, xreg(t, 0), xreg(t, 1), 5, 63)
	_, ok := in.(Sbfm)
	require.True(t, ok, "type = %T, want Sbfm", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"sbfm sp,rd",
			func() error {
				_, err := New().Sbfm(SP, xreg(t, 1), 5, 63)
				return err
			},
		},
		{
			"sbfm sp,rn",
			func() error {
				_, err := New().Sbfm(xreg(t, 0), SP, 5, 63)
				return err
			},
		},
		{
			"sbfm x,w widths",
			func() error {
				_, err := New().Sbfm(xreg(t, 0), wreg(t, 1), 5, 63)
				return err
			},
		},
		{
			"sbfm immr 64",
			func() error {
				_, err := New().Sbfm(xreg(t, 0), xreg(t, 1), 64, 63)
				return err
			},
		},
		{
			"sbfm imms 64",
			func() error {
				_, err := New().Sbfm(xreg(t, 0), xreg(t, 1), 5, 64)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorSbfm — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorSbfm(t *testing.T, rd, rn Reg, immr, imms uint32) Instr {
	t.Helper()
	in, err := New().Sbfm(rd, rn, immr, imms)
	require.NoError(t, err)
	return in
}
