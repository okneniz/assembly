package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdrpCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"adrp x0,#1",
			ctorAdrp(t, xreg(t, 0), 1),
			0xb0000000,
		},
		{
			"adrp x2,#0x10",
			ctorAdrp(t, xreg(t, 2), 0x10),
			0x90000082,
		},
		{
			"adrp x1,#-0x400",
			ctorAdrp(t, xreg(t, 1), -0x400),
			0x90ffe001,
		},
		{
			"adrp x0,#0xfffff",
			ctorAdrp(t, xreg(t, 0), 0xfffff),
			0xf07fffe0,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)

		// The absolute-page annotation of the decoded form is derived
		// from the instruction address and is not stored by the ctor —
		// compare the encoded fields, not the ObjDump text.
		back, ok := decodeOne(c.word, 0x1000).(Adrp)
		require.True(t, ok, "case %q: type = %T, want Adrp", c.name, back)
		want, wok := c.in.(Adrp)
		require.True(t, wok, "case %q: built %T, want Adrp", c.name, c.in)
		require.Equal(t, want.rd, back.rd, "case %q", c.name)
		require.Equal(t, want.off, back.off, "case %q", c.name)
	}

	in := ctorAdrp(t, xreg(t, 0), 1)
	_, ok := in.(Adrp)
	require.True(t, ok, "type = %T, want Adrp", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"adrp w rd",
			func() error {
				_, err := New().Adrp(wreg(t, 0), 1)
				return err
			},
		},
		{
			"adrp off 0x100000",
			func() error {
				_, err := New().Adrp(xreg(t, 0), 0x100000)
				return err
			},
		},
		{
			"adrp off -0x100001",
			func() error {
				_, err := New().Adrp(xreg(t, 0), -0x100001)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorAdrp — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorAdrp(t *testing.T, rd Reg, off int64) Instr {
	t.Helper()
	in, err := New().Adrp(rd, off)
	require.NoError(t, err)
	return in
}
