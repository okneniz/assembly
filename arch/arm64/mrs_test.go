package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestMrsCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"mrs x0,MIDR_EL1",
			ctorMrs(t, xreg(t, 0), "MIDR_EL1"),
			0xd5380000,
		},
		{
			"mrs x5,SCTLR_EL1",
			ctorMrs(t, xreg(t, 5), "SCTLR_EL1"),
			0xd5381005,
		},
		{
			"mrs x1,NZCV",
			ctorMrs(t, xreg(t, 1), "NZCV"),
			0xd53b4201,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorMrs(t, xreg(t, 0), "MIDR_EL1")
	_, ok := in.(Mrs)
	require.True(t, ok, "type = %T, want Mrs", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"mrs w0,rd",
			func() error {
				_, err := New().Mrs(wreg(t, 0), "MIDR_EL1")
				return err
			},
		},
		{
			"mrs bad sysreg",
			func() error {
				_, err := New().Mrs(xreg(t, 0), "NOT_A_SYSREG")
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorMrs — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorMrs(t *testing.T, rd Reg, sysreg string) Instr {
	t.Helper()
	in, err := New().Mrs(rd, sysreg)
	require.NoError(t, err)
	return in
}
