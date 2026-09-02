package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSubsImmCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"subs x0,x1,#0x42",
			ctorSubsImm(t, xreg(t, 0), xreg(t, 1), imm12(t, 0x42), NoSh12),
			0xf1010820,
		},
		{
			"subs xzr,x1,#1 (cmp)",
			ctorSubsImm(t, XZR, xreg(t, 1), imm12(t, 1), NoSh12),
			0xf100043f,
		},
		{
			"subs w2,w3,#7,lsl#12",
			ctorSubsImm(t, wreg(t, 2), wreg(t, 3), imm12(t, 7), LSL12),
			0x71401c62,
		},
		{
			"subs x2,sp,#0x10",
			ctorSubsImm(t, xreg(t, 2), SP, imm12(t, 0x10), NoSh12),
			0xf10043e2,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorSubsImm(t, xreg(t, 0), xreg(t, 1), imm12(t, 0x42), NoSh12)
	_, ok := in.(SubsImm)
	require.True(t, ok, "type = %T, want SubsImm", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"subs sp,rd",
			func() error {
				_, err := New().SubsImm(SP, xreg(t, 1), imm12(t, 1), NoSh12)
				return err
			},
		},
		{
			"subs xzr,rn",
			func() error {
				_, err := New().SubsImm(xreg(t, 0), XZR, imm12(t, 1), NoSh12)
				return err
			},
		},
		{
			"subs x,w widths",
			func() error {
				_, err := New().SubsImm(xreg(t, 0), wreg(t, 1), imm12(t, 1), NoSh12)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorSubsImm — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorSubsImm(t *testing.T, rd, rn Reg, imm Imm12, sh Sh12) Instr {
	t.Helper()
	in, err := New().SubsImm(rd, rn, imm, sh)
	require.NoError(t, err)
	return in
}
