package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClzBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"clz x1,x2",
			buildClz(t, xreg(t, 1), xreg(t, 2)),
			0xdac01041,
		},
		{
			"clz w1,w2",
			buildClz(t, wreg(t, 1), wreg(t, 2)),
			0x5ac01041,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
	}

	in := buildClz(t, xreg(t, 1), xreg(t, 2))
	_, ok := in.(Clz)
	require.True(t, ok, "type = %T, want Clz", in)
	for _, c := range []struct {
		name string
		call func() error
	}{
		{
			"clz x+w",
			func() error {
				_, err := New().Clz(xreg(t, 1), wreg(t, 2))
				return err
			},
		},
		{
			"clz sp",
			func() error {
				_, err := New().Clz(xreg(t, 1), SP)
				return err
			},
		},
		{
			"clz rd sp",
			func() error {
				_, err := New().Clz(SP, xreg(t, 2))
				return err
			},
		},
	} {
		assertErr(t, c.name, c.call())
	}
}

// buildClz — an instruction constructor wrapper for table literals:
// valid operands, an error is impossible by construction.
func buildClz(t *testing.T, rd, rn Reg) Instr {
	t.Helper()
	in, err := New().Clz(rd, rn)
	require.NoError(t, err)
	return in
}
