package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSturBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"stur x0,[x1]",
			buildStur(t, xreg(t, 0), xreg(t, 1), 0),
			0xf8000020,
		},
		{
			"stur w1,[x2,#8]",
			buildStur(t, wreg(t, 1), xreg(t, 2), 8),
			0xb8008041,
		},
		{
			"stur x3,[x4,#-256]",
			buildStur(t, xreg(t, 3), xreg(t, 4), -256),
			0xf8100083,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := buildStur(t, xreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Stur)
	require.True(t, ok, "type = %T, want Stur", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"stur sp,rt",
			func() error {
				_, err := New().Stur(SP, xreg(t, 1), 0)
				return err
			},
		},
		{
			"stur w1 base",
			func() error {
				_, err := New().Stur(xreg(t, 0), wreg(t, 1), 0)
				return err
			},
		},
		{
			"stur off -257",
			func() error {
				_, err := New().Stur(xreg(t, 0), xreg(t, 1), -257)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildStur — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildStur(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Stur(rt, rn, off)
	require.NoError(t, err)
	return in
}
