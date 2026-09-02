package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStlxrbCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"stlxrb w0,w1,[x2]",
			ctorStlxrb(t, wreg(t, 0), wreg(t, 1), xreg(t, 2)),
			0x0800fc41,
		},
		{
			"stlxrb wzr,w0,[x1]",
			ctorStlxrb(t, WZR, wreg(t, 0), xreg(t, 1)),
			0x081ffc20,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorStlxrb(t, wreg(t, 0), wreg(t, 1), xreg(t, 2))
	_, ok := in.(Stlxrb)
	require.True(t, ok, "type = %T, want Stlxrb", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"stlxrb x,rs",
			func() error {
				_, err := New().Stlxrb(xreg(t, 0), wreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"stlxrb x,rt",
			func() error {
				_, err := New().Stlxrb(wreg(t, 0), xreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"stlxrb w2 base",
			func() error {
				_, err := New().Stlxrb(wreg(t, 0), wreg(t, 1), wreg(t, 2))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorStlxrb — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorStlxrb(t *testing.T, rs, rt, rn Reg) Instr {
	t.Helper()
	in, err := New().Stlxrb(rs, rt, rn)
	require.NoError(t, err)
	return in
}
