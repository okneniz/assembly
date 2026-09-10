package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsrRegBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		word uint32
	}{
		{"asr x1,x2,x3", "x1", "x2", "x3", 0x9a032841},
		{"asr w1,w2,w3", "w1", "w2", "w3", 0x1a032841},
		{"asr x1,xzr,x3", "x1", "xzr", "x3", 0x9a032be1},
	}
	for _, c := range cases {
		in, err := New().AsrReg(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().AsrReg(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm))
	require.NoError(t, err)
	_, ok := in.(AsrReg)
	require.True(t, ok, "type = %T, want AsrReg", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
	}{
		{"asr x+w", "x0", "w1", "x2"},
		{"asr sp", "sp", "x1", "x2"},
		{"asr rn sp", "x0", "sp", "x2"},
		{"asr rm sp", "x0", "x1", "sp"},
	}
	for _, c := range errCases {
		_, err := New().AsrReg(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm))
		assertErr(t, c.name, err)
	}
}
