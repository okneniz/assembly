package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRev16Build(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		word uint32
	}{
		{"rev16 x0,x1", "x0", "x1", 0xdac00420},
		{"rev16 w2,w3", "w2", "w3", 0x5ac00462},
		{"rev16 x1,xzr", "x1", "xzr", 0xdac007e1},
	}
	for _, c := range cases {
		in, err := New().Rev16(reg(t, c.rd), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Rev16(reg(t, first.rd), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Rev16)
	require.True(t, ok, "type = %T, want Rev16", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
	}{
		{"rev16 sp,rd", "sp", "x1"},
		{"rev16 sp,rn", "x0", "sp"},
		{"rev16 x,w widths", "x0", "w1"},
	}
	for _, c := range errCases {
		_, err := New().Rev16(reg(t, c.rd), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
