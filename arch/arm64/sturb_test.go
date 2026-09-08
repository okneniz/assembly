package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSturbBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"sturb w0,[x1]",
			buildSturb(t, wreg(t, 0), xreg(t, 1), 0),
			0x38000020,
		},
		{
			"sturb wzr,[x29,#255]",
			buildSturb(t, WZR, xreg(t, 29), 255),
			0x380ff3bf,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := buildSturb(t, wreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Sturb)
	require.True(t, ok, "type = %T, want Sturb", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"sturb x,rt",
			func() error {
				_, err := New().Sturb(xreg(t, 0), xreg(t, 1), 0)
				return err
			},
		},
		{
			"sturb w1 base",
			func() error {
				_, err := New().Sturb(wreg(t, 0), wreg(t, 1), 0)
				return err
			},
		},
		{
			"sturb off 256",
			func() error {
				_, err := New().Sturb(wreg(t, 0), xreg(t, 1), 256)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildSturb — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildSturb(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Sturb(rt, rn, off)
	require.NoError(t, err)
	return in
}
