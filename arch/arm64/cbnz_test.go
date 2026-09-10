package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCbnzBuild(t *testing.T) {
	cases := []struct {
		name   string
		rt     string
		target int64
		word   uint32
	}{
		{"cbnz x2,0x1008", "x2", 8, 0xb5000042},
		{"cbnz w0,0x1014", "w0", 20, 0x350000a0},
		{"cbnz wzr,0x1000", "wzr", 0, 0x3500001f},
		{"cbnz x0,0x100ffc", "x0", 1048572, 0xb57fffe0},
	}
	for _, c := range cases {
		in, err := New().Cbnz(reg(t, c.rt), c.target)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Cbnz(reg(t, first.rt), first.target)
	require.NoError(t, err)
	_, ok := in.(Cbnz)
	require.True(t, ok, "type = %T, want Cbnz", in)

	errCases := []struct {
		name   string
		rt     string
		target int64
	}{
		{"cbnz sp,rt", "sp", 0x1000},
	}
	for _, c := range errCases {
		_, err := New().Cbnz(reg(t, c.rt), c.target)
		assertErr(t, c.name, err)
	}
}
