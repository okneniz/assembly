package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEorImmCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"eor x0,x1,#0x7",
			ctorEorImm(t, xreg(t, 0), xreg(t, 1), 0x7),
			0xd2400820,
		},
		{
			"eor x2,x3,#0xffff",
			ctorEorImm(t, xreg(t, 2), xreg(t, 3), 0xffff),
			0xd2403c62,
		},
		{
			"eor w0,w1,#0x00ff00ff",
			ctorEorImm(t, wreg(t, 0), wreg(t, 1), 0x00ff00ff),
			0x52009c20,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorEorImm(t, xreg(t, 0), xreg(t, 1), 0x7)
	_, ok := in.(EorImm)
	require.True(t, ok, "type = %T, want EorImm", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"eor imm sp",
			func() error {
				_, err := New().EorImm(SP, xreg(t, 1), 0x7)
				return err
			},
		},
		{
			"eor imm rn sp",
			func() error {
				_, err := New().EorImm(xreg(t, 0), SP, 0x7)
				return err
			},
		},
		{
			"eor imm x+w",
			func() error {
				_, err := New().EorImm(xreg(t, 0), wreg(t, 1), 0x7)
				return err
			},
		},
		{
			"eor imm not encodable",
			func() error {
				_, err := New().EorImm(xreg(t, 0), xreg(t, 1), 0x55)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorEorImm — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorEorImm(t *testing.T, rd, rn Reg, imm uint64) Instr {
	t.Helper()
	in, err := New().EorImm(rd, rn, imm)
	require.NoError(t, err)
	return in
}
