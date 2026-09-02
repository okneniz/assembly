package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCsinvCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"csinv x5,x6,x7,al",
			ctorCsinv(t, xreg(t, 5), xreg(t, 6), xreg(t, 7), "al"),
			0xda87e0c5,
		},
		{
			"csinv w2,w3,wzr,le",
			ctorCsinv(t, wreg(t, 2), wreg(t, 3), WZR, "le"),
			0x5a9fd062,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorCsinv(t, xreg(t, 5), xreg(t, 6), xreg(t, 7), "al")
	_, ok := in.(Csinv)
	require.True(t, ok, "type = %T, want Csinv", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"csinv sp,rd",
			func() error {
				_, err := New().Csinv(SP, xreg(t, 6), xreg(t, 7), "al")
				return err
			},
		},
		{
			"csinv w,x widths",
			func() error {
				_, err := New().Csinv(xreg(t, 5), wreg(t, 6), xreg(t, 7), "al")
				return err
			},
		},
		{
			"csinv bad cond",
			func() error {
				_, err := New().Csinv(xreg(t, 5), xreg(t, 6), xreg(t, 7), "foo")
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorCsinv — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorCsinv(t *testing.T, rd, rn, rm Reg, cond string) Instr {
	t.Helper()
	in, err := New().Csinv(rd, rn, rm, cond)
	require.NoError(t, err)
	return in
}
