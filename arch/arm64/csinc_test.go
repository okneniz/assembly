package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCsincCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"csinc x1,x2,x3,eq",
			ctorCsinc(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), "eq"),
			0x9a830441,
		},
		{
			"csinc w0,wzr,w1,ne",
			ctorCsinc(t, wreg(t, 0), WZR, wreg(t, 1), "ne"),
			0x1a8117e0,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorCsinc(t, xreg(t, 1), xreg(t, 2), xreg(t, 3), "eq")
	_, ok := in.(Csinc)
	require.True(t, ok, "type = %T, want Csinc", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"csinc sp,rd",
			func() error {
				_, err := New().Csinc(SP, xreg(t, 2), xreg(t, 3), "eq")
				return err
			},
		},
		{
			"csinc w,x widths",
			func() error {
				_, err := New().Csinc(xreg(t, 1), wreg(t, 2), xreg(t, 3), "eq")
				return err
			},
		},
		{
			"csinc bad cond",
			func() error {
				_, err := New().Csinc(xreg(t, 1), xreg(t, 2), xreg(t, 3), "foo")
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorCsinc — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorCsinc(t *testing.T, rd, rn, rm Reg, cond string) Instr {
	t.Helper()
	in, err := New().Csinc(rd, rn, rm, cond)
	require.NoError(t, err)
	return in
}
