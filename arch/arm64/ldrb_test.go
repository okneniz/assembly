package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdrbCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldrb w0,[x1]",
			ctorLdrb(t, wreg(t, 0), xreg(t, 1), 0),
			0x39400020,
		},
		{
			"ldrb w2,[sp,#1]",
			ctorLdrb(t, wreg(t, 2), SP, 1),
			0x394007e2,
		},
		{
			"ldrb wzr,[x2,#0xfff]",
			ctorLdrb(t, WZR, xreg(t, 2), 0xfff),
			0x397ffc5f,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorLdrb(t, wreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Ldrb)
	require.True(t, ok, "type = %T, want Ldrb", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldrb x,rt",
			func() error {
				_, err := New().Ldrb(xreg(t, 0), xreg(t, 1), 0)
				return err
			},
		},
		{
			"ldrb w1 base",
			func() error {
				_, err := New().Ldrb(wreg(t, 0), wreg(t, 1), 0)
				return err
			},
		},
		{
			"ldrb off -1",
			func() error {
				_, err := New().Ldrb(wreg(t, 0), xreg(t, 1), -1)
				return err
			},
		},
		{
			"ldrb off 0x1000",
			func() error {
				_, err := New().Ldrb(wreg(t, 0), xreg(t, 1), 0x1000)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorLdrb — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorLdrb(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Ldrb(rt, rn, off)
	require.NoError(t, err)
	return in
}
