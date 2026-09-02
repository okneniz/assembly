package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsrRegCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"asr x1,x2,x3",
			ctorAsrReg(t, xreg(t, 1), xreg(t, 2), xreg(t, 3)),
			0x9a032841,
		},
		{
			"asr w1,w2,w3",
			ctorAsrReg(t, wreg(t, 1), wreg(t, 2), wreg(t, 3)),
			0x1a032841,
		},
		{
			"asr x1,xzr,x3",
			ctorAsrReg(t, xreg(t, 1), XZR, xreg(t, 3)),
			0x9a032be1,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorAsrReg(t, xreg(t, 1), xreg(t, 2), xreg(t, 3))
	_, ok := in.(AsrReg)
	require.True(t, ok, "type = %T, want AsrReg", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"asr x+w",
			func() error {
				_, err := New().AsrReg(xreg(t, 0), wreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"asr sp",
			func() error {
				_, err := New().AsrReg(SP, xreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"asr rn sp",
			func() error {
				_, err := New().AsrReg(xreg(t, 0), SP, xreg(t, 2))
				return err
			},
		},
		{
			"asr rm sp",
			func() error {
				_, err := New().AsrReg(xreg(t, 0), xreg(t, 1), SP)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorAsrReg — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorAsrReg(t *testing.T, rd, rn, rm Reg) Instr {
	t.Helper()
	in, err := New().AsrReg(rd, rn, rm)
	require.NoError(t, err)
	return in
}
