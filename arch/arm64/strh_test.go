package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStrhCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"strh w0,[x1]",
			ctorStrh(t, wreg(t, 0), xreg(t, 1), 0),
			0x79000020,
		},
		{
			"strh w2,[x3,#0x1ffe]",
			ctorStrh(t, wreg(t, 2), xreg(t, 3), 0x1ffe),
			0x793ffc62,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorStrh(t, wreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Strh)
	require.True(t, ok, "type = %T, want Strh", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"strh x,rt",
			func() error {
				_, err := New().Strh(xreg(t, 0), xreg(t, 1), 0)
				return err
			},
		},
		{
			"strh w1 base",
			func() error {
				_, err := New().Strh(wreg(t, 0), wreg(t, 1), 0)
				return err
			},
		},
		{
			"strh off 1",
			func() error {
				_, err := New().Strh(wreg(t, 0), xreg(t, 1), 1)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorStrh — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorStrh(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Strh(rt, rn, off)
	require.NoError(t, err)
	return in
}
