package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestUdivCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"udiv x0,x1,x2",
			ctorUdiv(t, xreg(t, 0), xreg(t, 1), xreg(t, 2)),
			0x9ac20820,
		},
		{
			"udiv w4,w5,w6",
			ctorUdiv(t, wreg(t, 4), wreg(t, 5), wreg(t, 6)),
			0x1ac608a4,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorUdiv(t, xreg(t, 0), xreg(t, 1), xreg(t, 2))
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

// ctorUdiv — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorUdiv(t *testing.T, rd, rn, rm Reg) Instr {
	t.Helper()
	in, err := New().Udiv(rd, rn, rm)
	require.NoError(t, err)
	return in
}
