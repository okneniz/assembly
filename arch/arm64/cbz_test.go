package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCbzBuild(t *testing.T) {
	cases := []struct {
		name   string
		rt     string
		target int64
		word   uint32
	}{
		{"cbz x0,0x1000", "x0", 0, 0xb4000000},
		{"cbz w3,0x1014", "w3", 20, 0x340000a3},
		{"cbz x1,0xfe8", "x1", -24, 0xb4ffff41},
		{"cbz x0,0x800", "x0", -2048, 0xb4ffc000},
		{"cbz x0,0x100ffc", "x0", 1048572, 0xb47fffe0},
	}
	for _, c := range cases {
		in, err := New().Cbz(reg(t, c.rt), c.target)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Cbz(reg(t, first.rt), first.target)
	require.NoError(t, err)
	_, ok := in.(Cbz)
	require.True(t, ok, "type = %T, want Cbz", in)

	errCases := []struct {
		name   string
		rt     string
		target int64
	}{
		{"cbz sp,rt", "sp", 0x1000},
	}
	for _, c := range errCases {
		_, err := New().Cbz(reg(t, c.rt), c.target)
		assertErr(t, c.name, err)
	}
}
