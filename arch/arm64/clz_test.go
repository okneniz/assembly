package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClzBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		word uint32
	}{
		{"clz x1,x2", "x1", "x2", 0xdac01041},
		{"clz w1,w2", "w1", "w2", 0x5ac01041},
	}
	for _, c := range cases {
		in, err := New().Clz(reg(t, c.rd), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Clz(reg(t, first.rd), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Clz)
	require.True(t, ok, "type = %T, want Clz", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
	}{
		{"clz x+w", "x1", "w2"},
		{"clz sp", "x1", "sp"},
		{"clz rd sp", "sp", "x2"},
	}
	for _, c := range errCases {
		_, err := New().Clz(reg(t, c.rd), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
