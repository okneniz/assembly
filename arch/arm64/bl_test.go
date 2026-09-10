package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBlBuild(t *testing.T) {
	cases := []struct {
		name string
		off  int64
		word uint32
	}{
		{"bl 0x1000", 0, 0x94000000},
		{"bl 0x1240", 576, 0x94000090},
		{"bl 0x0", -4096, 0x97fffc00},
		{"bl 0x8000ffc", 134217724, 0x95ffffff},
	}
	for _, c := range cases {
		in := New().Bl(c.off)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
		back, derr := decodeOne(c.word)
		_ = derr
		require.Equal(t, in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}
}
