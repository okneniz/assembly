package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAdrCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"adr x0,#0x10",
			ctorAdr(t, xreg(t, 0), 0x10),
			0x10000080,
		},
		{
			"adr x1,#-0x1000",
			ctorAdr(t, xreg(t, 1), -0x1000),
			0x10ff8001,
		},
		{
			"adr x2,#0x7ff",
			ctorAdr(t, xreg(t, 2), 0x7ff),
			0x70003fe2,
		},
		{
			"adr x3,#0xfffff",
			ctorAdr(t, xreg(t, 3), 0xfffff),
			0x707fffe3,
		},
		{
			"adr x4,#-0x100000",
			ctorAdr(t, xreg(t, 4), -0x100000),
			0x10800004,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorAdr(t, xreg(t, 0), 0x10)
	_, ok := in.(Adr)
	require.True(t, ok, "type = %T, want Adr", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"adr w rd",
			func() error {
				_, err := New().Adr(wreg(t, 0), 0x10)
				return err
			},
		},
		{
			"adr off 0x100000",
			func() error {
				_, err := New().Adr(xreg(t, 0), 0x100000)
				return err
			},
		},
		{
			"adr off -0x100001",
			func() error {
				_, err := New().Adr(xreg(t, 0), -0x100001)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorAdr — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorAdr(t *testing.T, rd Reg, off int64) Instr {
	t.Helper()
	in, err := New().Adr(rd, off)
	require.NoError(t, err)
	return in
}
