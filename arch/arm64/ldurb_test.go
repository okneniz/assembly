package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdurbCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"ldurb w0,[x1]",
			ctorLdurb(t, wreg(t, 0), xreg(t, 1), 0),
			0x38400020,
		},
		{
			"ldurb w3,[x4,#-256]",
			ctorLdurb(t, wreg(t, 3), xreg(t, 4), -256),
			0x38500083,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word, 0x1000)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := ctorLdurb(t, wreg(t, 0), xreg(t, 1), 0)
	_, ok := in.(Ldurb)
	require.True(t, ok, "type = %T, want Ldurb", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"ldurb x,rt",
			func() error {
				_, err := New().Ldurb(xreg(t, 0), xreg(t, 1), 0)
				return err
			},
		},
		{
			"ldurb w1 base",
			func() error {
				_, err := New().Ldurb(wreg(t, 0), wreg(t, 1), 0)
				return err
			},
		},
		{
			"ldurb off 256",
			func() error {
				_, err := New().Ldurb(wreg(t, 0), xreg(t, 1), 256)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorLdurb — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorLdurb(t *testing.T, rt, rn Reg, off Off) Instr {
	t.Helper()
	in, err := New().Ldurb(rt, rn, off)
	require.NoError(t, err)
	return in
}
