package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSturhBuild(t *testing.T) {
	cases := []struct {
		name string
		rt   string
		rn   string
		off  Off
		word uint32
	}{
		{"sturh w0,[x1]", "w0", "x1", 0, 0x78000020},
		{"sturh w7,[x8,#-9]", "w7", "x8", -9, 0x781f7107},
	}
	for _, c := range cases {
		in, err := New().Sturh(reg(t, c.rt), reg(t, c.rn), c.off)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Sturh(reg(t, first.rt), reg(t, first.rn), first.off)
	require.NoError(t, err)
	_, ok := in.(Sturh)
	require.True(t, ok, "type = %T, want Sturh", in)

	errCases := []struct {
		name string
		rt   string
		rn   string
		off  Off
	}{
		{"sturh x,rt", "x0", "x1", 0},
		{"sturh w1 base", "w0", "w1", 0},
		{"sturh off -257", "w0", "x1", -257},
	}
	for _, c := range errCases {
		_, err := New().Sturh(reg(t, c.rt), reg(t, c.rn), c.off)
		assertErr(t, c.name, err)
	}
}
