package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBBuild(t *testing.T) {
	cases := []struct {
		name string
		off  int64
		word uint32
	}{
		{"b 0x1008", 8, 0x14000002},
		{"b 0x1240", 576, 0x14000090},
		{"b 0x0", -4096, 0x17fffc00},
	}
	for _, c := range cases {
		in := New().B(c.off)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
		back, derr := decodeOne(c.word)
		_ = derr
		require.Equal(t, in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}
}
