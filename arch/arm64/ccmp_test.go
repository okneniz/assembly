package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCcmpCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ccmp x1,x2,#3,eq",
			ctorCcmp(t, xreg(t, 1), xreg(t, 2), 3, "eq"),
			0xfa420023,
		},
		{
			"ccmp x0,xzr,#0xf,ne",
			ctorCcmp(t, xreg(t, 0), XZR, 0xf, "ne"),
			0xfa5f100f,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorCcmp(t, xreg(t, 1), xreg(t, 2), 3, "eq")
	_, ok := in.(Ccmp)
	require.True(t, ok, "type = %T, want Ccmp", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ccmp w form",
			func() error {
				_, err := New().Ccmp(wreg(t, 1), xreg(t, 2), 3, "eq")
				return err
			},
		},
		{
			"ccmp rm w form",
			func() error {
				_, err := New().Ccmp(xreg(t, 1), wreg(t, 2), 3, "eq")
				return err
			},
		},
		{
			"ccmp nzcv=0x10",
			func() error {
				_, err := New().Ccmp(xreg(t, 1), xreg(t, 2), 0x10, "eq")
				return err
			},
		},
		{
			"ccmp bad cond",
			func() error {
				_, err := New().Ccmp(xreg(t, 1), xreg(t, 2), 3, "foo")
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorCcmp — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorCcmp(t *testing.T, rn, rm Reg, nzcv uint32, cond string) Instr {
	t.Helper()
	in, err := New().Ccmp(rn, rm, nzcv, cond)
	require.NoError(t, err)
	return in
}
