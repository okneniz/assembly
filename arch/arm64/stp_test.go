package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStpCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"stp x0,x1,[x2]",
			ctorStp(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), 0),
			0xa9000440,
		},
		{
			"stp x29,x30,[sp,#-16]",
			ctorStp(t, xreg(t, 29), xreg(t, 30), SP, -16),
			0xa93f7bfd,
		},
		{
			"stp w0,w1,[x2,#252]",
			ctorStp(t, wreg(t, 0), wreg(t, 1), xreg(t, 2), 252),
			0x291f8440,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorStp(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), 0)
	_, ok := in.(Stp)
	require.True(t, ok, "type = %T, want Stp", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"stp sp,rt",
			func() error {
				_, err := New().Stp(SP, xreg(t, 1), xreg(t, 2), 0)
				return err
			},
		},
		{
			"stp sp,rt2",
			func() error {
				_, err := New().Stp(xreg(t, 0), SP, xreg(t, 2), 0)
				return err
			},
		},
		{
			"stp w2 base",
			func() error {
				_, err := New().Stp(xreg(t, 0), xreg(t, 1), wreg(t, 2), 0)
				return err
			},
		},
		{
			"stp x,w widths",
			func() error {
				_, err := New().Stp(xreg(t, 0), wreg(t, 1), xreg(t, 2), 0)
				return err
			},
		},
		{
			"stp off 4 in x form",
			func() error {
				_, err := New().Stp(xreg(t, 0), xreg(t, 1), xreg(t, 2), 4)
				return err
			},
		},
		{
			"stp off -520 in x form",
			func() error {
				_, err := New().Stp(xreg(t, 0), xreg(t, 1), xreg(t, 2), -520)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorStp — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorStp(t *testing.T, rt, rt2, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Stp(rt, rt2, rn, off)
	require.NoError(t, err)
	return in
}
