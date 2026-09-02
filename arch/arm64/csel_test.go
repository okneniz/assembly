package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCselCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"csel x0,x1,x2,eq",
			ctorCsel(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), "eq"),
			0x9a820020,
		},
		{
			"csel w3,w4,w5,hs",
			ctorCsel(t, wreg(t, 3), wreg(t, 4), wreg(t, 5), "hs"),
			0x1a852083,
		},
		{
			"csel x0,xzr,x1,gt",
			ctorCsel(t, xreg(t, 0), XZR, xreg(t, 1), "gt"),
			0x9a81c3e0,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorCsel(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), "eq")
	_, ok := in.(Csel)
	require.True(t, ok, "type = %T, want Csel", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"csel sp,rd",
			func() error {
				_, err := New().Csel(SP, xreg(t, 1), xreg(t, 2), "eq")
				return err
			},
		},
		{
			"csel sp,rn",
			func() error {
				_, err := New().Csel(xreg(t, 0), SP, xreg(t, 2), "eq")
				return err
			},
		},
		{
			"csel sp,rm",
			func() error {
				_, err := New().Csel(xreg(t, 0), xreg(t, 1), SP, "eq")
				return err
			},
		},
		{
			"csel w,x widths",
			func() error {
				_, err := New().Csel(xreg(t, 0), wreg(t, 1), xreg(t, 2), "eq")
				return err
			},
		},
		{
			"csel bad cond",
			func() error {
				_, err := New().Csel(xreg(t, 0), xreg(t, 1), xreg(t, 2), "foo")
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorCsel — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorCsel(t *testing.T, rd, rn, rm Reg, cond string) Instr {
	t.Helper()
	in, err := New().Csel(rd, rn, rm, cond)
	require.NoError(t, err)
	return in
}
