package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBlrBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"blr x1",
			buildBlr(t, xreg(t, 1)),
			0xd63f0020,
		},
		{
			"blr x30",
			buildBlr(t, xreg(t, 30)),
			0xd63f03c0,
		},
		{
			"blr xzr",
			buildBlr(t, XZR),
			0xd63f03e0,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := buildBlr(t, xreg(t, 1))
	_, ok := in.(Blr)
	require.True(t, ok, "type = %T, want Blr", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"blr w1",
			func() error {
				_, err := New().Blr(wreg(t, 1))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildBlr — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildBlr(t *testing.T, rn Reg) Instr {
	t.Helper()
	in, err := New().Blr(rn)
	require.NoError(t, err)
	return in
}
