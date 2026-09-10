package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSmulhBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		word uint32
	}{
		{"smulh x0,x1,x2", "x0", "x1", "x2", 0x9b427c20},
		{"smulh xzr,x1,x2", "xzr", "x1", "x2", 0x9b427c3f},
	}
	for _, c := range cases {
		in, err := New().Smulh(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Smulh(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm))
	require.NoError(t, err)
	_, ok := in.(Smulh)
	require.True(t, ok, "type = %T, want Smulh", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
	}{
		{"smulh w form", "w0", "w1", "w2"},
		{"smulh w,rn", "x0", "w1", "x2"},
		{"smulh w,rm", "x0", "x1", "w2"},
	}
	for _, c := range errCases {
		_, err := New().Smulh(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		assertErr(t, c.name, err)
	}
}
