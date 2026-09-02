package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddExtCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"add x1,x2,x3,uxtb#0",
			ctorAddExt(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), "uxtb", 0),
			0x8b230041,
		},
		{
			"add x0,x1,x2,uxtx#7",
			ctorAddExt(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), "uxtx", 7),
			0x8b227c20,
		},
		{
			"add w1,w2,w3,uxtw#2",
			ctorAddExt(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), "uxtw", 2),
			0x0b234841,
		},
		{
			"add sp,sp,x3,sxtx#3",
			ctorAddExt(t, SP, SP, xreg(t, 3), "sxtx", 3),
			0x8b23efff,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorAddExt(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), "uxtb", 0)
	_, ok := in.(AddExt)
	require.True(t, ok, "type = %T, want AddExt", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"add ext rd xzr",
			func() error {
				_, err := New().AddExt(XZR, xreg(t, 1), xreg(t, 2), "uxtb", 0)
				return err
			},
		},
		{
			"add ext rn xzr",
			func() error {
				_, err := New().AddExt(xreg(t, 0), XZR, xreg(t, 2), "uxtb", 0)
				return err
			},
		},
		{
			"add ext rm xzr",
			func() error {
				_, err := New().AddExt(xreg(t, 0), xreg(t, 1), XZR, "uxtb", 0)
				return err
			},
		},
		{
			"add ext x+w",
			func() error {
				_, err := New().AddExt(xreg(t, 0), wreg(t, 1), xreg(t, 2), "uxtb", 0)
				return err
			},
		},
		{
			"add ext bad ext",
			func() error {
				_, err := New().AddExt(xreg(t, 0), xreg(t, 1), xreg(t, 2), "uxtx2", 0)
				return err
			},
		},
		{
			"add ext imm3=8",
			func() error {
				_, err := New().AddExt(xreg(t, 0), xreg(t, 1), xreg(t, 2), "uxtx", 8)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorAddExt — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorAddExt(t *testing.T, rd, rn, rm Reg, ext string, imm3 uint32) Instr {
	t.Helper()
	in, err := New().AddExt(rd, rn, rm, ext, imm3)
	require.NoError(t, err)
	return in
}
