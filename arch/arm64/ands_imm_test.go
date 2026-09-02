package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAndsImmCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"tst x1,#0x7",
			ctorAndsImm(t, XZR, xreg(t, 1), 0x7),
			0xf240083f,
		},
		{
			"ands x2,x3,#0xffff0000ffff0000",
			ctorAndsImm(t, xreg(t, 2), xreg(t, 3), 0xffff0000ffff0000),
			0xf2103c62,
		},
		{
			"ands w2,w3,#0x00ff00ff",
			ctorAndsImm(t, wreg(t, 2), wreg(t, 3), 0x00ff00ff),
			0x72009c62,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorAndsImm(t, xreg(t, 2), xreg(t, 3), 0x7)
	_, ok := in.(AndsImm)
	require.True(t, ok, "type = %T, want AndsImm", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ands imm sp",
			func() error {
				_, err := New().AndsImm(SP, xreg(t, 1), 0x7)
				return err
			},
		},
		{
			"ands imm rn sp",
			func() error {
				_, err := New().AndsImm(xreg(t, 0), SP, 0x7)
				return err
			},
		},
		{
			"ands imm x+w",
			func() error {
				_, err := New().AndsImm(xreg(t, 0), wreg(t, 1), 0x7)
				return err
			},
		},
		{
			"ands imm not encodable",
			func() error {
				_, err := New().AndsImm(xreg(t, 0), xreg(t, 1), ^uint64(0))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorAndsImm — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorAndsImm(t *testing.T, rd, rn Reg, imm uint64) Instr {
	t.Helper()
	in, err := New().AndsImm(rd, rn, imm)
	require.NoError(t, err)
	return in
}
