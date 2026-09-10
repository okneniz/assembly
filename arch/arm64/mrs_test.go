package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMrsBuild(t *testing.T) {
	cases := []struct {
		name   string
		rd     string
		sysreg string
		word   uint32
	}{
		{"mrs x0,MIDR_EL1", "x0", "MIDR_EL1", 0xd5380000},
		{"mrs x5,SCTLR_EL1", "x5", "SCTLR_EL1", 0xd5381005},
		{"mrs x1,NZCV", "x1", "NZCV", 0xd53b4201},
	}
	for _, c := range cases {
		in, err := New().Mrs(reg(t, c.rd), c.sysreg)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Mrs(reg(t, first.rd), first.sysreg)
	require.NoError(t, err)
	_, ok := in.(Mrs)
	require.True(t, ok, "type = %T, want Mrs", in)

	errCases := []struct {
		name   string
		rd     string
		sysreg string
	}{
		{"mrs w0,rd", "w0", "MIDR_EL1"},
		{"mrs bad sysreg", "x0", "NOT_A_SYSREG"},
	}
	for _, c := range errCases {
		_, err := New().Mrs(reg(t, c.rd), c.sysreg)
		assertErr(t, c.name, err)
	}
}
