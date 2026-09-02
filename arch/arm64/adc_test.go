package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdcCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"adc x0,x1,x2",
			ctorAdc(t, xreg(t, 0), xreg(t, 1), xreg(t, 2)),
			0x9a020020,
		},
		{
			"adc x5,xzr,x7",
			ctorAdc(t, xreg(t, 5), XZR, xreg(t, 7)),
			0x9a0703e5,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorAdc(t, xreg(t, 0), xreg(t, 1), xreg(t, 2))
	_, ok := in.(Adc)
	require.True(t, ok, "type = %T, want Adc", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"adc w form",
			func() error {
				_, err := New().Adc(wreg(t, 0), xreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"adc sp",
			func() error {
				_, err := New().Adc(SP, xreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"adc rn w form",
			func() error {
				_, err := New().Adc(xreg(t, 0), wreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"adc rm sp",
			func() error {
				_, err := New().Adc(xreg(t, 0), xreg(t, 1), SP)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorAdc — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorAdc(t *testing.T, rd, rn, rm Reg) Instr {
	t.Helper()
	in, err := New().Adc(rd, rn, rm)
	require.NoError(t, err)
	return in
}
