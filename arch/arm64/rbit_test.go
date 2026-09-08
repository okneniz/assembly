package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRbitBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"rbit x0,x1",
			buildRbit(t, xreg(t, 0), xreg(t, 1)),
			0xdac00020,
		},
		{
			"rbit w2,w3",
			buildRbit(t, wreg(t, 2), wreg(t, 3)),
			0x5ac00062,
		},
		{
			"rbit xzr,x30",
			buildRbit(t, XZR, xreg(t, 30)),
			0xdac003df,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := buildRbit(t, xreg(t, 0), xreg(t, 1))
	_, ok := in.(Rbit)
	require.True(t, ok, "type = %T, want Rbit", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"rbit sp,rd",
			func() error {
				_, err := New().Rbit(SP, xreg(t, 1))
				return err
			},
		},
		{
			"rbit sp,rn",
			func() error {
				_, err := New().Rbit(xreg(t, 0), SP)
				return err
			},
		},
		{
			"rbit x,w widths",
			func() error {
				_, err := New().Rbit(xreg(t, 0), wreg(t, 1))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildRbit — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildRbit(t *testing.T, rd, rn Reg) Instr {
	t.Helper()
	in, err := New().Rbit(rd, rn)
	require.NoError(t, err)
	return in
}
