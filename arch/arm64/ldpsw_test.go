package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdpswBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldpsw x0,x1,[x2,#4]",
			buildLdpsw(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), 4),
			0x69408440,
		},
		{
			"ldpsw x2,x3,[x4,#-256]",
			buildLdpsw(t, xreg(t, 2), xreg(t, 3), xreg(t, 4), -256),
			0x69600c82,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := buildLdpsw(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), 4)
	_, ok := in.(Ldpsw)
	require.True(t, ok, "type = %T, want Ldpsw", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldpsw w,rt",
			func() error {
				_, err := New().Ldpsw(wreg(t, 0), xreg(t, 1), xreg(t, 2), 4)
				return err
			},
		},
		{
			"ldpsw w,rt2",
			func() error {
				_, err := New().Ldpsw(xreg(t, 0), wreg(t, 1), xreg(t, 2), 4)
				return err
			},
		},
		{
			"ldpsw w2 base",
			func() error {
				_, err := New().Ldpsw(xreg(t, 0), xreg(t, 1), wreg(t, 2), 4)
				return err
			},
		},
		{
			"ldpsw off 2",
			func() error {
				_, err := New().Ldpsw(xreg(t, 0), xreg(t, 1), xreg(t, 2), 2)
				return err
			},
		},
		{
			"ldpsw off 260",
			func() error {
				_, err := New().Ldpsw(xreg(t, 0), xreg(t, 1), xreg(t, 2), 260)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildLdpsw — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildLdpsw(t *testing.T, rt, rt2, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Ldpsw(rt, rt2, rn, off)
	require.NoError(t, err)
	return in
}
