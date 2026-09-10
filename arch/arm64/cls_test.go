package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClsBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		word uint32
	}{
		{"cls x1,x2", "x1", "x2", 0xdac01441},
		{"cls w1,w2", "w1", "w2", 0x5ac01441},
	}
	for _, c := range cases {
		in, err := New().Cls(reg(t, c.rd), reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Cls(reg(t, first.rd), reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Cls)
	require.True(t, ok, "type = %T, want Cls", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
	}{
		{"cls x+w", "x1", "w2"},
		{"cls sp", "sp", "x2"},
		{"cls rn sp", "x1", "sp"},
	}
	for _, c := range errCases {
		_, err := New().Cls(reg(t, c.rd), reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
