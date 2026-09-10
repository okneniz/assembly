package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsrBuild(t *testing.T) {
	cases := []struct {
		name   string
		sysreg string
		rt     string
		word   uint32
	}{
		{"msr SCTLR_EL1,x0", "SCTLR_EL1", "x0", 0xd5181000},
		{"msr NZCV,x1", "NZCV", "x1", 0xd51b4201},
		{"msr DAIF,x2", "DAIF", "x2", 0xd51b4222},
	}
	for _, c := range cases {
		in, err := New().Msr(c.sysreg, reg(t, c.rt))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Msr(first.sysreg, reg(t, first.rt))
	require.NoError(t, err)
	_, ok := in.(Msr)
	require.True(t, ok, "type = %T, want Msr", in)

	errCases := []struct {
		name   string
		sysreg string
		rt     string
	}{
		{"msr w0,rt", "SCTLR_EL1", "w0"},
		{"msr bad sysreg", "NOT_A_SYSREG", "x0"},
	}
	for _, c := range errCases {
		_, err := New().Msr(c.sysreg, reg(t, c.rt))
		assertErr(t, c.name, err)
	}
}
