package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrrImmBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"orr x0,x1,#0x7",
			buildOrrImm(t, xreg(t, 0), xreg(t, 1), 0x7),
			0xb2400820,
		},
		{
			"orr x0,xzr,#0xffff",
			buildOrrImm(t, xreg(t, 0), XZR, 0xffff),
			0xb2403fe0,
		},
		{
			"orr w0,w1,#0x00ff00ff",
			buildOrrImm(t, wreg(t, 0), wreg(t, 1), 0x00ff00ff),
			0x32009c20,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := buildOrrImm(t, xreg(t, 0), xreg(t, 1), 0x7)
	_, ok := in.(OrrImm)
	require.True(t, ok, "type = %T, want OrrImm", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"orr imm sp",
			func() error {
				_, err := New().OrrImm(SP, xreg(t, 1), 0x7)
				return err
			},
		},
		{
			"orr imm rn sp",
			func() error {
				_, err := New().OrrImm(xreg(t, 0), SP, 0x7)
				return err
			},
		},
		{
			"orr imm x+w",
			func() error {
				_, err := New().OrrImm(xreg(t, 0), wreg(t, 1), 0x7)
				return err
			},
		},
		{
			"orr imm not encodable",
			func() error {
				_, err := New().OrrImm(xreg(t, 0), xreg(t, 1), 0)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildOrrImm — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildOrrImm(t *testing.T, rd, rn Reg, imm uint64) Instr {
	t.Helper()
	in, err := New().OrrImm(rd, rn, imm)
	require.NoError(t, err)
	return in
}
