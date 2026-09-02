package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsubCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"msub x1,x2,x3,x4",
			ctorMsub(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), xreg(t, 4)),
			0x9b039041,
		},
		{
			"mneg w1,w2,w3",
			ctorMsub(t, wreg(t, 1), wreg(t, 2), wreg(t, 3), WZR),
			0x1b03fc41,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorMsub(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), xreg(t, 4))
	_, ok := in.(Msub)
	require.True(t, ok, "type = %T, want Msub", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"msub x+w",
			func() error {
				_, err := New().Msub(xreg(t, 1), wreg(t, 2), xreg(t, 3), xreg(t, 4))
				return err
			},
		},
		{
			"msub sp",
			func() error {
				_, err := New().Msub(xreg(t, 1), xreg(t, 2), xreg(t, 3), SP)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorMsub — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorMsub(t *testing.T, rd, rn, rm, ra Reg) Instr {
	t.Helper()
	in, err := New().Msub(rd, rn, rm, ra)
	require.NoError(t, err)
	return in
}
