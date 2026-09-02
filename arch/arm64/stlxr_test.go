package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStlxrCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"stlxr w0,x1,[x2]",
			ctorStlxr(t, wreg(t, 0), xreg(t, 1), xreg(t, 2)),
			0xc800fc41,
		},
		{
			"stlxr w1,w2,[x3]",
			ctorStlxr(t, wreg(t, 1), wreg(t, 2), xreg(t, 3)),
			0x8801fc62,
		},
		{
			"stlxr wzr,x0,[sp]",
			ctorStlxr(t, WZR, xreg(t, 0), SP),
			0xc81fffe0,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorStlxr(t, wreg(t, 0), xreg(t, 1), xreg(t, 2))
	_, ok := in.(Stlxr)
	require.True(t, ok, "type = %T, want Stlxr", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"stlxr x,rs",
			func() error {
				_, err := New().Stlxr(xreg(t, 0), xreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"stlxr sp,rt",
			func() error {
				_, err := New().Stlxr(wreg(t, 0), SP, xreg(t, 2))
				return err
			},
		},
		{
			"stlxr w2 base",
			func() error {
				_, err := New().Stlxr(wreg(t, 0), xreg(t, 1), wreg(t, 2))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorStlxr — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorStlxr(t *testing.T, rs, rt, rn Reg) Instr {
	t.Helper()
	in, err := New().Stlxr(rs, rt, rn)
	require.NoError(t, err)
	return in
}
