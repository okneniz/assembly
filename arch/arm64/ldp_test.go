package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdpCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldp x0,x1,[x2]",
			ctorLdp(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), 0),
			0xa9400440,
		},
		{
			"ldp x29,x30,[sp]",
			ctorLdp(t, xreg(t, 29), xreg(t, 30), SP, 0),
			0xa9407bfd,
		},
		{
			"ldp w3,w4,[x5,#8]",
			ctorLdp(t, wreg(t, 3), wreg(t, 4), xreg(t, 5), 8),
			0x294110a3,
		},
		{
			"ldp x0,x1,[x2,#-512]",
			ctorLdp(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), -512),
			0xa9600440,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorLdp(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), 0)
	_, ok := in.(Ldp)
	require.True(t, ok, "type = %T, want Ldp", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldp sp,rt",
			func() error {
				_, err := New().Ldp(SP, xreg(t, 1), xreg(t, 2), 0)
				return err
			},
		},
		{
			"ldp sp,rt2",
			func() error {
				_, err := New().Ldp(xreg(t, 0), SP, xreg(t, 2), 0)
				return err
			},
		},
		{
			"ldp w2 base",
			func() error {
				_, err := New().Ldp(xreg(t, 0), xreg(t, 1), wreg(t, 2), 0)
				return err
			},
		},
		{
			"ldp x,w widths",
			func() error {
				_, err := New().Ldp(xreg(t, 0), wreg(t, 1), xreg(t, 2), 0)
				return err
			},
		},
		{
			"ldp off 4 in x form",
			func() error {
				_, err := New().Ldp(xreg(t, 0), xreg(t, 1), xreg(t, 2), 4)
				return err
			},
		},
		{
			"ldp off 520 in x form",
			func() error {
				_, err := New().Ldp(xreg(t, 0), xreg(t, 1), xreg(t, 2), 520)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorLdp — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorLdp(t *testing.T, rt, rt2, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Ldp(rt, rt2, rn, off)
	require.NoError(t, err)
	return in
}
