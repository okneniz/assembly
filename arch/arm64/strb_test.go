package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStrbBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"strb w0,[x1]",
			buildStrb(t, wreg(t, 0), xreg(t, 1), 0),
			0x39000020,
		},
		{
			"strb wzr,[sp,#1]",
			buildStrb(t, WZR, SP, 1),
			0x390007ff,
		},
		{
			"strb w2,[x3,#0xfff]",
			buildStrb(t, wreg(t, 2), xreg(t, 3), 0xfff),
			0x393ffc62,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := buildStrb(t, wreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Strb)
	require.True(t, ok, "type = %T, want Strb", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"strb x,rt",
			func() error {
				_, err := New().Strb(xreg(t, 0), xreg(t, 1), 0)
				return err
			},
		},
		{
			"strb w1 base",
			func() error {
				_, err := New().Strb(wreg(t, 0), wreg(t, 1), 0)
				return err
			},
		},
		{
			"strb off -1",
			func() error {
				_, err := New().Strb(wreg(t, 0), xreg(t, 1), -1)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildStrb — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildStrb(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Strb(rt, rn, off)
	require.NoError(t, err)
	return in
}
