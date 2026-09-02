package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdrhCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldrh w0,[x1,#2]",
			ctorLdrh(t, wreg(t, 0), xreg(t, 1), 2),
			0x79400420,
		},
		{
			"ldrh w1,[sp,#0x1ffe]",
			ctorLdrh(t, wreg(t, 1), SP, 0x1ffe),
			0x797fffe1,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorLdrh(t, wreg(t, 0), xreg(t, 1), 2)
	_, ok := in.(Ldrh)
	require.True(t, ok, "type = %T, want Ldrh", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldrh x,rt",
			func() error {
				_, err := New().Ldrh(xreg(t, 0), xreg(t, 1), 2)
				return err
			},
		},
		{
			"ldrh w1 base",
			func() error {
				_, err := New().Ldrh(wreg(t, 0), wreg(t, 1), 2)
				return err
			},
		},
		{
			"ldrh off 1",
			func() error {
				_, err := New().Ldrh(wreg(t, 0), xreg(t, 1), 1)
				return err
			},
		},
		{
			"ldrh off 0x2000",
			func() error {
				_, err := New().Ldrh(wreg(t, 0), xreg(t, 1), 0x2000)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorLdrh — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorLdrh(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Ldrh(rt, rn, off)
	require.NoError(t, err)
	return in
}
