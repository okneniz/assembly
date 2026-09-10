package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdrpBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		off  int64
		word uint32
	}{
		{"adrp x0,#1", "x0", 1, 0xb0000000},
		{"adrp x2,#0x10", "x2", 0x10, 0x90000082},
		{"adrp x1,#-0x400", "x1", -0x400, 0x90ffe001},
		{"adrp x0,#0xfffff", "x0", 0xfffff, 0xf07fffe0},
	}
	for _, c := range cases {
		in, err := New().Adrp(reg(t, c.rd), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Adrp(reg(t, first.rd), first.off)
	require.NoError(t, err)
	_, ok := in.(Adrp)
	require.True(t, ok, "type = %T, want Adrp", in)

	errCases := []struct {
		name string
		rd   string
		off  int64
	}{
		{"adrp w rd", "w0", 1},
		{"adrp off 0x100000", "x0", 0x100000},
		{"adrp off -0x100001", "x0", -0x100001},
	}
	for _, c := range errCases {
		_, err := New().Adrp(reg(t, c.rd), c.off)
		assertErr(t, c.name, err)
	}
}
