package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBCtor(t *testing.T) {
	cases := []struct {
		name string
		in   Instr
		word uint32
	}{
		{
			"b 0x1008",
			New().B(0x1008),
			0x14000002,
		},
		{
			"b 0x1240",
			New().B(0x1240),
			0x14000090,
		},
		{
			"b 0x0",
			New().B(0),
			0x17fffc00,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.in), "case %q", c.name)
			back := decodeOne(c.word, 0x1000)
			require.Equal(t, c.in.ObjDump(disasm.DefaultViewCtx()),
				back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
		})
	}
}
