package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSubsExtCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"subs x0,x1,x2,sxtx#2",
			ctorSubsExt(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), "sxtx", 2),
			0xeb22e820,
		},
		{
			"subs xzr,x1,x2,sxtx (cmp)",
			ctorSubsExt(t, XZR, xreg(t, 1), xreg(t, 2), "sxtx", 0),
			0xeb22e03f,
		},
		{
			"subs w1,w2,w3,uxth",
			ctorSubsExt(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), "uxth", 0),
			0x6b232041,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorSubsExt(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), "sxtx", 2)
	_, ok := in.(SubsExt)
	require.True(t, ok, "type = %T, want SubsExt", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"subs sp,rd",
			func() error {
				_, err := New().SubsExt(SP, xreg(t, 1), xreg(t, 2), "sxtx", 2)
				return err
			},
		},
		{
			"subs xzr,rn",
			func() error {
				_, err := New().SubsExt(xreg(t, 0), XZR, xreg(t, 2), "sxtx", 2)
				return err
			},
		},
		{
			"subs xzr,rm",
			func() error {
				_, err := New().SubsExt(xreg(t, 0), xreg(t, 1), XZR, "sxtx", 2)
				return err
			},
		},
		{
			"subs x,w widths",
			func() error {
				_, err := New().SubsExt(xreg(t, 0), wreg(t, 1), xreg(t, 2), "sxtx", 2)
				return err
			},
		},
		{
			"subs bad ext",
			func() error {
				_, err := New().SubsExt(xreg(t, 0), xreg(t, 1), xreg(t, 2), "foo", 2)
				return err
			},
		},
		{
			"subs imm3 8",
			func() error {
				_, err := New().SubsExt(xreg(t, 0), xreg(t, 1), xreg(t, 2), "sxtx", 8)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorSubsExt — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorSubsExt(t *testing.T, rd, rn, rm Reg, ext string, imm3 uint32) Instr {
	t.Helper()
	in, err := New().SubsExt(rd, rn, rm, ext, imm3)
	require.NoError(t, err)
	return in
}
