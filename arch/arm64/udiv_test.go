package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestUdivBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"udiv x0,x1,x2",
			buildUdiv(t, xreg(t, 0), xreg(t, 1), xreg(t, 2)),
			0x9ac20820,
		},
		{
			"udiv w4,w5,w6",
			buildUdiv(t, wreg(t, 4), wreg(t, 5), wreg(t, 6)),
			0x1ac608a4,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := buildUdiv(t, xreg(t, 0), xreg(t, 1), xreg(t, 2))
	_, ok := in.(Udiv)
	require.True(t, ok, "type = %T, want Udiv", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"udiv sp,rd",
			func() error {
				_, err := New().Udiv(SP, xreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"udiv sp,rn",
			func() error {
				_, err := New().Udiv(xreg(t, 0), SP, xreg(t, 2))
				return err
			},
		},
		{
			"udiv sp,rm",
			func() error {
				_, err := New().Udiv(xreg(t, 0), xreg(t, 1), SP)
				return err
			},
		},
		{
			"udiv x,w widths",
			func() error {
				_, err := New().Udiv(xreg(t, 0), wreg(t, 1), xreg(t, 2))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildUdiv — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildUdiv(t *testing.T, rd, rn, rm Reg) Instr {
	t.Helper()
	in, err := New().Udiv(rd, rn, rm)
	require.NoError(t, err)
	return in
}
