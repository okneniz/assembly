package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStlxrBuild(t *testing.T) {
	cases := []struct {
		name string
		rs   string
		rt   string
		rn   string
		word uint32
	}{
		{"stlxr w0,x1,[x2]", "w0", "x1", "x2", 0xc800fc41},
		{"stlxr w1,w2,[x3]", "w1", "w2", "x3", 0x8801fc62},
		{"stlxr wzr,x0,[sp]", "wzr", "x0", "sp", 0xc81fffe0},
	}
	for _, c := range cases {
		in, err := New().Stlxr(reg(t, c.rs), reg(t, c.rt), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Stlxr(reg(t, first.rs), reg(t, first.rt), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Stlxr)
	require.True(t, ok, "type = %T, want Stlxr", in)

	errCases := []struct {
		name string
		rs   string
		rt   string
		rn   string
	}{
		{"stlxr x,rs", "x0", "x1", "x2"},
		{"stlxr sp,rt", "w0", "sp", "x2"},
		{"stlxr w2 base", "w0", "x1", "w2"},
	}
	for _, c := range errCases {
		_, err := New().Stlxr(reg(t, c.rs), reg(t, c.rt), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
