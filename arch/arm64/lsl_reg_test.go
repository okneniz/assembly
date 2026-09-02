package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLslRegCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"lsl x1,x2,x3",
			ctorLslReg(t, xreg(t, 1), xreg(t, 2), xreg(t, 3)),
			0x9a032041,
		},
		{
			"lsl w1,w2,w3",
			ctorLslReg(t, wreg(t, 1), wreg(t, 2), wreg(t, 3)),
			0x1a032041,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorLslReg(t, xreg(t, 1), xreg(t, 2), xreg(t, 3))
	_, ok := in.(LslReg)
	require.True(t, ok, "type = %T, want LslReg", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"lsl x+w",
			func() error {
				_, err := New().LslReg(xreg(t, 0), wreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"lsl sp",
			func() error {
				_, err := New().LslReg(xreg(t, 0), xreg(t, 1), SP)
				return err
			},
		},
		{
			"lsl rd sp",
			func() error {
				_, err := New().LslReg(SP, xreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"lsl rn sp",
			func() error {
				_, err := New().LslReg(xreg(t, 0), SP, xreg(t, 2))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorLslReg — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorLslReg(t *testing.T, rd, rn, rm Reg) Instr {
	t.Helper()
	in, err := New().LslReg(rd, rn, rm)
	require.NoError(t, err)
	return in
}
