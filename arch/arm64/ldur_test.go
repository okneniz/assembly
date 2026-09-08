package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdurBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldur x0,[x1]",
			buildLdur(t, xreg(t, 0), xreg(t, 1), 0),
			0xf8400020,
		},
		{
			"ldur x1,[x2,#-1]",
			buildLdur(t, xreg(t, 1), xreg(t, 2), -1),
			0xf85ff041,
		},
		{
			"ldur w2,[sp,#255]",
			buildLdur(t, wreg(t, 2), SP, 255),
			0xb84ff3e2,
		},
		{
			"ldur w3,[x4,#-256]",
			buildLdur(t, wreg(t, 3), xreg(t, 4), -256),
			0xb8500083,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := buildLdur(t, xreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Ldur)
	require.True(t, ok, "type = %T, want Ldur", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldur sp,rt",
			func() error {
				_, err := New().Ldur(SP, xreg(t, 1), 0)
				return err
			},
		},
		{
			"ldur w1 base",
			func() error {
				_, err := New().Ldur(xreg(t, 0), wreg(t, 1), 0)
				return err
			},
		},
		{
			"ldur off 256",
			func() error {
				_, err := New().Ldur(xreg(t, 0), xreg(t, 1), 256)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildLdur — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildLdur(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Ldur(rt, rn, off)
	require.NoError(t, err)
	return in
}
