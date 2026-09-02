package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClsCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"cls x1,x2",
			ctorCls(t, xreg(t, 1), xreg(t, 2)),
			0xdac01441,
		},
		{
			"cls w1,w2",
			ctorCls(t, wreg(t, 1), wreg(t, 2)),
			0x5ac01441,
		},
	}
	for _, c := range cases {
		got := ctorWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := ctorCls(t, xreg(t, 1), xreg(t, 2))
	_, ok := in.(Cls)
	require.True(t, ok, "type = %T, want Cls", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"cls x+w",
			func() error {
				_, err := New().Cls(xreg(t, 1), wreg(t, 2))
				return err
			},
		},
		{
			"cls sp",
			func() error {
				_, err := New().Cls(SP, xreg(t, 2))
				return err
			},
		},
		{
			"cls rn sp",
			func() error {
				_, err := New().Cls(xreg(t, 1), SP)
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// ctorCls — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func ctorCls(t *testing.T, rd, rn Reg) Instr {
	t.Helper()
	in, err := New().Cls(rd, rn)
	require.NoError(t, err)
	return in
}
