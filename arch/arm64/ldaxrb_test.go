package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdaxrbBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldaxrb w0,[x1]",
			buildLdaxrb(t, wreg(t, 0), xreg(t, 1)),
			0x085ffc20,
		},
		{
			"ldaxrb w5,[sp]",
			buildLdaxrb(t, wreg(t, 5), SP),
			0x085fffe5,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := buildLdaxrb(t, wreg(t, 0), xreg(t, 1))
	_, ok := in.(Ldaxrb)
	require.True(t, ok, "type = %T, want Ldaxrb", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldaxrb x,rt",
			func() error {
				_, err := New().Ldaxrb(xreg(t, 0), xreg(t, 1))
				return err
			},
		},
		{
			"ldaxrb w1 base",
			func() error {
				_, err := New().Ldaxrb(wreg(t, 0), wreg(t, 1))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildLdaxrb — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildLdaxrb(t *testing.T, rt, rn Reg) Instr {
	t.Helper()
	in, err := New().Ldaxrb(rt, rn)
	require.NoError(t, err)
	return in
}
