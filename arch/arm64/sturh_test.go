package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSturhCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"sturh w0,[x1]",
			ctorSturh(t, wreg(t, 0), xreg(t, 1), 0),
			0x78000020,
		},
		{
			"sturh w7,[x8,#-9]",
			ctorSturh(t, wreg(t, 7), xreg(t, 8), -9),
			0x781f7107,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorSturh(t, wreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Sturh)
	require.True(t, ok, "type = %T, want Sturh", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"sturh x,rt",
			func() error {
				_, err := New().Sturh(xreg(t, 0), xreg(t, 1), 0)
				return err
			},
		},
		{
			"sturh w1 base",
			func() error {
				_, err := New().Sturh(wreg(t, 0), wreg(t, 1), 0)
				return err
			},
		},
		{
			"sturh off -257",
			func() error {
				_, err := New().Sturh(wreg(t, 0), xreg(t, 1), -257)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorSturh — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorSturh(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Sturh(rt, rn, off)
	require.NoError(t, err)
	return in
}
