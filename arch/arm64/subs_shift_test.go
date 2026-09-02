package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSubsShiftCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"subs x0,x1,x2",
			ctorSubsShift(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), imm6(t, 0), LSL),
			0xeb020020,
		},
		{
			"subs x0,x1,x2,lsl#3",
			ctorSubsShift(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), imm6(t, 3), LSL),
			0xeb020c20,
		},
		{
			"subs xzr,x1,x2,asr#4 (cmp)",
			ctorSubsShift(t, XZR, xreg(t, 1), xreg(t, 2), imm6(t, 4), ASR),
			0xeb82103f,
		},
		{
			"subs w3,w4,w5,lsl#2",
			ctorSubsShift(t, wreg(t, 3), wreg(t, 4), wreg(t, 5), imm6(t, 2), LSL),
			0x6b050883,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorSubsShift(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), imm6(t, 0), LSL)
	_, ok := in.(SubsShift)
	require.True(t, ok, "type = %T, want SubsShift", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"subs sp,rd",
			func() error {
				_, err := New().SubsShift(SP, xreg(t, 1), xreg(t, 2), imm6(t, 0), LSL)
				return err
			},
		},
		{
			"subs sp,rn",
			func() error {
				_, err := New().SubsShift(xreg(t, 0), SP, xreg(t, 2), imm6(t, 0), LSL)
				return err
			},
		},
		{
			"subs sp,rm",
			func() error {
				_, err := New().SubsShift(xreg(t, 0), xreg(t, 1), SP, imm6(t, 0), LSL)
				return err
			},
		},
		{
			"subs x,w widths",
			func() error {
				_, err := New().SubsShift(xreg(t, 0), wreg(t, 1), xreg(t, 2), imm6(t, 0), LSL)
				return err
			},
		},
		{
			"subs ror shift",
			func() error {
				_, err := New().SubsShift(xreg(t, 0), xreg(t, 1), xreg(t, 2), imm6(t, 0), ROR)
				return err
			},
		},
		{
			"subs w,#32 shift",
			func() error {
				_, err := New().SubsShift(wreg(t, 0), wreg(t, 1), wreg(t, 2), imm6(t, 32), LSL)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorSubsShift — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorSubsShift(t *testing.T, rd, rn, rm Reg, imm Imm6, sh Shift) Instr {
	t.Helper()
	in, err := New().SubsShift(rd, rn, rm, imm, sh)
	require.NoError(t, err)
	return in
}
