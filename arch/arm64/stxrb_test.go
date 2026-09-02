package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStxrbCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"stxrb w0,w1,[x2]",
			ctorStxrb(t, wreg(t, 0), wreg(t, 1), xreg(t, 2)),
			0x08000041,
		},
		{
			"stxrb wzr,w1,[sp]",
			ctorStxrb(t, WZR, wreg(t, 1), SP),
			0x081f03e1,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorStxrb(t, wreg(t, 0), wreg(t, 1), xreg(t, 2))
	_, ok := in.(Stxrb)
	require.True(t, ok, "type = %T, want Stxrb", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"stxrb x,rs",
			func() error {
				_, err := New().Stxrb(xreg(t, 0), wreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"stxrb x,rt",
			func() error {
				_, err := New().Stxrb(wreg(t, 0), xreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"stxrb w2 base",
			func() error {
				_, err := New().Stxrb(wreg(t, 0), wreg(t, 1), wreg(t, 2))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorStxrb — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorStxrb(t *testing.T, rs, rt, rn Reg) Instr {
	t.Helper()
	in, err := New().Stxrb(rs, rt, rn)
	require.NoError(t, err)
	return in
}
