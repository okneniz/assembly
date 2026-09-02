package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddsImmCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"cmn x1,#0x42",
			ctorAddsImm(t, XZR, xreg(t, 1), imm12(t, 0x42), NoSh12),
			0xb101083f,
		},
		{
			"adds x1,x2,#0x10,lsl#12",
			ctorAddsImm(t, xreg(t, 1), xreg(t, 2), imm12(t, 0x10), LSL12),
			0xb1404041,
		},
		{
			"adds w2,w3,#7",
			ctorAddsImm(t, wreg(t, 2), wreg(t, 3), imm12(t, 7), NoSh12),
			0x31001c62,
		},
		{
			"adds x2,x3,#0xfff",
			ctorAddsImm(t, xreg(t, 2), xreg(t, 3), imm12(t, 0xfff), NoSh12),
			0xb13ffc62,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorAddsImm(t, xreg(t, 1), xreg(t, 2), imm12(t, 1), NoSh12)
	_, ok := in.(AddsImm)
	require.True(t, ok, "type = %T, want AddsImm", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"adds rd sp",
			func() error {
				_, err := New().AddsImm(SP, xreg(t, 1), imm12(t, 1), NoSh12)
				return err
			},
		},
		{
			"adds rn xzr",
			func() error {
				_, err := New().AddsImm(xreg(t, 0), XZR, imm12(t, 1), NoSh12)
				return err
			},
		},
		{
			"adds x+w",
			func() error {
				_, err := New().AddsImm(xreg(t, 0), wreg(t, 1), imm12(t, 1), NoSh12)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorAddsImm — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorAddsImm(t *testing.T, rd, rn Reg, imm Imm12, sh Sh12) Instr {
	t.Helper()
	in, err := New().AddsImm(rd, rn, imm, sh)
	require.NoError(t, err)
	return in
}
