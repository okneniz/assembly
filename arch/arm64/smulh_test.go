package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSmulhCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"smulh x0,x1,x2",
			ctorSmulh(t, xreg(t, 0), xreg(t, 1), xreg(t, 2)),
			0x9b427c20,
		},
		{
			"smulh xzr,x1,x2",
			ctorSmulh(t, XZR, xreg(t, 1), xreg(t, 2)),
			0x9b427c3f,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorSmulh(t, xreg(t, 0), xreg(t, 1), xreg(t, 2))
	_, ok := in.(Smulh)
	require.True(t, ok, "type = %T, want Smulh", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"smulh w form",
			func() error {
				_, err := New().Smulh(wreg(t, 0), wreg(t, 1), wreg(t, 2))
				return err
			},
		},
		{
			"smulh w,rn",
			func() error {
				_, err := New().Smulh(xreg(t, 0), wreg(t, 1), xreg(t, 2))
				return err
			},
		},
		{
			"smulh w,rm",
			func() error {
				_, err := New().Smulh(xreg(t, 0), xreg(t, 1), wreg(t, 2))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorSmulh — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorSmulh(t *testing.T, rd, rn, rm Reg) Instr {
	t.Helper()
	in, err := New().Smulh(rd, rn, rm)
	require.NoError(t, err)
	return in
}
