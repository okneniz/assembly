package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddsExtCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"cmn x1,x2,uxtb#0",
			ctorAddsExt(t, XZR, xreg(t, 1), xreg(t, 2), "uxtb", 0),
			0xab22003f,
		},
		{
			"adds w1,w2,w3,sxtw#1",
			ctorAddsExt(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), "sxtw", 1),
			0x2b23c441,
		},
		{
			"adds x1,sp,x3,uxtx#0",
			ctorAddsExt(t, xreg(t, 1), SP, xreg(t, 3), "uxtx", 0),
			0xab2363e1,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorAddsExt(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), "uxtb", 0)
	_, ok := in.(AddsExt)
	require.True(t, ok, "type = %T, want AddsExt", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"adds ext rd sp",
			func() error {
				_, err := New().AddsExt(SP, xreg(t, 1), xreg(t, 2), "uxtb", 0)
				return err
			},
		},
		{
			"adds ext rn xzr",
			func() error {
				_, err := New().AddsExt(xreg(t, 0), XZR, xreg(t, 2), "uxtb", 0)
				return err
			},
		},
		{
			"adds ext rm xzr",
			func() error {
				_, err := New().AddsExt(xreg(t, 0), xreg(t, 1), XZR, "uxtb", 0)
				return err
			},
		},
		{
			"adds ext x+w",
			func() error {
				_, err := New().AddsExt(xreg(t, 0), wreg(t, 1), xreg(t, 2), "uxtb", 0)
				return err
			},
		},
		{
			"adds ext bad ext",
			func() error {
				_, err := New().AddsExt(xreg(t, 0), xreg(t, 1), xreg(t, 2), "sxtb2", 0)
				return err
			},
		},
		{
			"adds ext imm3=8",
			func() error {
				_, err := New().AddsExt(xreg(t, 0), xreg(t, 1), xreg(t, 2), "uxtx", 8)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorAddsExt — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorAddsExt(t *testing.T, rd, rn, rm Reg, ext string, imm3 uint32) Instr {
	t.Helper()
	in, err := New().AddsExt(rd, rn, rm, ext, imm3)
	require.NoError(t, err)
	return in
}
