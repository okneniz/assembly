package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAndImmCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"and x0,x1,#0x7",
			ctorAndImm(t, xreg(t, 0), xreg(t, 1), 0x7),
			0x92400820,
		},
		{
			"and x0,x1,#0xffff",
			ctorAndImm(t, xreg(t, 0), xreg(t, 1), 0xffff),
			0x92403c20,
		},
		{
			"and w0,w1,#0x00ff00ff",
			ctorAndImm(t, wreg(t, 0), wreg(t, 1), 0x00ff00ff),
			0x12009c20,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorAndImm(t, xreg(t, 0), xreg(t, 1), 0x7)
	_, ok := in.(AndImm)
	require.True(t, ok, "type = %T, want AndImm", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"and imm sp",
			func() error {
				_, err := New().AndImm(SP, xreg(t, 1), 0x7)
				return err
			},
		},
		{
			"and imm rn sp",
			func() error {
				_, err := New().AndImm(xreg(t, 0), SP, 0x7)
				return err
			},
		},
		{
			"and imm x+w",
			func() error {
				_, err := New().AndImm(xreg(t, 0), wreg(t, 1), 0x7)
				return err
			},
		},
		{
			"and imm not encodable",
			func() error {
				_, err := New().AndImm(xreg(t, 0), xreg(t, 1), 0)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorAndImm — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorAndImm(t *testing.T, rd, rn Reg, imm uint64) Instr {
	t.Helper()
	in, err := New().AndImm(rd, rn, imm)
	require.NoError(t, err)
	return in
}
