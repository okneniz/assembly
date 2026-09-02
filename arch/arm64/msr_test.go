package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestMsrCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"msr SCTLR_EL1,x0",
			ctorMsr(t, "SCTLR_EL1", xreg(t, 0)),
			0xd5181000,
		},
		{
			"msr NZCV,x1",
			ctorMsr(t, "NZCV", xreg(t, 1)),
			0xd51b4201,
		},
		{
			"msr DAIF,x2",
			ctorMsr(t, "DAIF", xreg(t, 2)),
			0xd51b4222,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorMsr(t, "SCTLR_EL1", xreg(t, 0))
	_, ok := in.(Msr)
	require.True(t, ok, "type = %T, want Msr", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"msr w0,rt",
			func() error {
				_, err := New().Msr("SCTLR_EL1", wreg(t, 0))
				return err
			},
		},
		{
			"msr bad sysreg",
			func() error {
				_, err := New().Msr("NOT_A_SYSREG", xreg(t, 0))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorMsr — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorMsr(t *testing.T, sysreg string, rt Reg) Instr {
	t.Helper()
	in, err := New().Msr(sysreg, rt)
	require.NoError(t, err)
	return in
}
