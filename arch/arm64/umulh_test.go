package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestUmulhBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"umulh x3,x4,x5",
			buildUmulh(t, xreg(t, 3), xreg(t, 4), xreg(t, 5)),
			0x9bc57c83,
		},
		{
			"umulh x0,xzr,x2",
			buildUmulh(t, xreg(t, 0), XZR, xreg(t, 2)),
			0x9bc27fe0,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := buildUmulh(t, xreg(t, 3), xreg(t, 4), xreg(t, 5))
	_, ok := in.(Umulh)
	require.True(t, ok, "type = %T, want Umulh", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"umulh w form",
			func() error {
				_, err := New().Umulh(wreg(t, 3), wreg(t, 4), wreg(t, 5))
				return err
			},
		},
		{
			"umulh w,rn",
			func() error {
				_, err := New().Umulh(xreg(t, 3), wreg(t, 4), xreg(t, 5))
				return err
			},
		},
		{
			"umulh w,rm",
			func() error {
				_, err := New().Umulh(xreg(t, 3), xreg(t, 4), wreg(t, 5))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildUmulh — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildUmulh(t *testing.T, rd, rn, rm Reg) Instr {
	t.Helper()
	in, err := New().Umulh(rd, rn, rm)
	require.NoError(t, err)
	return in
}
