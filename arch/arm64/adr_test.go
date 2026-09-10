package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdrBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		off  int64
		word uint32
	}{
		{"adr x0,#0x10", "x0", 0x10, 0x10000080},
		{"adr x1,#-0x1000", "x1", -0x1000, 0x10ff8001},
		{"adr x2,#0x7ff", "x2", 0x7ff, 0x70003fe2},
		{"adr x3,#0xfffff", "x3", 0xfffff, 0x707fffe3},
		{"adr x4,#-0x100000", "x4", -0x100000, 0x10800004},
	}
	for _, c := range cases {
		in, err := New().Adr(reg(t, c.rd), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Adr(reg(t, first.rd), first.off)
	require.NoError(t, err)
	_, ok := in.(Adr)
	require.True(t, ok, "type = %T, want Adr", in)

	errCases := []struct {
		name string
		rd   string
		off  int64
	}{
		{"adr w rd", "w0", 0x10},
		{"adr off 0x100000", "x0", 0x100000},
		{"adr off -0x100001", "x0", -0x100001},
	}
	for _, c := range errCases {
		_, err := New().Adr(reg(t, c.rd), c.off)
		assertErr(t, c.name, err)
	}
}
