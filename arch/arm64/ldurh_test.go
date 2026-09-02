package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdurhCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldurh w0,[x1]",
			ctorLdurh(t, wreg(t, 0), xreg(t, 1), 0),
			0x78400020,
		},
		{
			"ldurh w5,[x6,#-3]",
			ctorLdurh(t, wreg(t, 5), xreg(t, 6), -3),
			0x785fd0c5,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorLdurh(t, wreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Ldurh)
	require.True(t, ok, "type = %T, want Ldurh", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldurh x,rt",
			func() error {
				_, err := New().Ldurh(xreg(t, 0), xreg(t, 1), 0)
				return err
			},
		},
		{
			"ldurh w1 base",
			func() error {
				_, err := New().Ldurh(wreg(t, 0), wreg(t, 1), 0)
				return err
			},
		},
		{
			"ldurh off 256",
			func() error {
				_, err := New().Ldurh(wreg(t, 0), xreg(t, 1), 256)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorLdurh — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorLdurh(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Ldurh(rt, rn, off)
	require.NoError(t, err)
	return in
}
