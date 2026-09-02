package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRorRegCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ror x0,x1,x2",
			ctorRorReg(t, xreg(t, 0), xreg(t, 1), xreg(t, 2)),
			0x9a022c20,
		},
		{
			"ror x3,xzr,x4",
			ctorRorReg(t, xreg(t, 3), XZR, xreg(t, 4)),
			0x9a042fe3,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorRorReg(t, xreg(t, 0), xreg(t, 1), xreg(t, 2))
	_, ok := in.(RorReg)
	require.True(t, ok, "type = %T, want RorReg", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ror w form",
			func() error {
				_, err := New().RorReg(wreg(t, 0), wreg(t, 1), wreg(t, 2))
				return err
			},
		},
		{
			"ror w,rn",
			func() error {
				_, err := New().RorReg(xreg(t, 0), wreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"ror w,rm",
			func() error {
				_, err := New().RorReg(xreg(t, 0), xreg(t, 1), wreg(t, 2))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorRorReg — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorRorReg(t *testing.T, rd, rn, rm Reg) Instr {
	t.Helper()
	in, err := New().RorReg(rd, rn, rm)
	require.NoError(t, err)
	return in
}
