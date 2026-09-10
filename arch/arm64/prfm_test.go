package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrfmBuild(t *testing.T) {
	cases := []struct {
		name string
		rn   string
		word uint32
	}{
		{"prfm pldl1keep,[x1]", "x1", 0xf9800020},
		{"prfm pldl1keep,[sp]", "sp", 0xf98003e0},
	}
	for _, c := range cases {
		in, err := New().Prfm(reg(t, c.rn))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Prfm(reg(t, first.rn))
	require.NoError(t, err)
	_, ok := in.(Prfm)
	require.True(t, ok, "type = %T, want Prfm", in)

	errCases := []struct {
		name string
		rn   string
	}{
		{"prfm w1 base", "w1"},
	}
	for _, c := range errCases {
		_, err := New().Prfm(reg(t, c.rn))
		assertErr(t, c.name, err)
	}
}
