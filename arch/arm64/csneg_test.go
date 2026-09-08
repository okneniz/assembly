package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCsnegBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"csneg x0,x1,x2,mi",
			buildCsneg(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), "mi"),
			0xda824420,
		},
		{
			"csneg w4,wzr,wzr,eq",
			buildCsneg(t, wreg(t, 4), WZR, WZR, "eq"),
			0x5a9f07e4,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := buildCsneg(t, xreg(t, 0), xreg(t, 1), xreg(t, 2), "mi")
	_, ok := in.(Csneg)
	require.True(t, ok, "type = %T, want Csneg", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"csneg sp,rd",
			func() error {
				_, err := New().Csneg(SP, xreg(t, 1), xreg(t, 2), "mi")
				return err
			},
		},
		{
			"csneg w,x widths",
			func() error {
				_, err := New().Csneg(xreg(t, 0), wreg(t, 1), xreg(t, 2), "mi")
				return err
			},
		},
		{
			"csneg bad cond",
			func() error {
				_, err := New().Csneg(xreg(t, 0), xreg(t, 1), xreg(t, 2), "foo")
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildCsneg — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildCsneg(t *testing.T, rd, rn, rm Reg, cond string) Instr {
	t.Helper()
	in, err := New().Csneg(rd, rn, rm, cond)
	require.NoError(t, err)
	return in
}
