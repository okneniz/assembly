package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBlBuild(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"bl 0x1000",
			New().Bl(0),
			0x94000000,
		},
		{
			"bl 0x1240",
			New().Bl(576),
			0x94000090,
		},
		{
			"bl 0x0",
			New().Bl(-4096),
			0x97fffc00,
		},
		{
			"bl 0x8000ffc",
			New().Bl(134217724),
			0x95ffffff,
		},
	}
	for _, c := range cases {
		got := buildWord(t, c.in)
		require.Equal(t, c.word, got, "case %q", c.name)
		back := decodeOne(c.word)
		require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	in := New().Bl(0)
	_, ok := in.(Bl)
	require.True(t, ok, "type = %T, want Bl", in)
}
