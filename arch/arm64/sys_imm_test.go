package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSvcBrkBuild(t *testing.T) {
	cases := []struct {
		name string
		imm  int64
		word uint32
	}{
		{"svc #0x80", 0x80, 0xd4001001},
	}
	for _, c := range cases {
		in := New().Svc(imm16(t, c.imm))
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
		back, derr := decodeOne(c.word)
		_ = derr
		require.Equal(t, in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	brkCases := []struct {
		name string
		imm  int64
		word uint32
	}{
		{"brk #1", 1, 0xd4200020},
		{"brk #0", 0, 0xd4200000},
	}
	for _, c := range brkCases {
		in := New().Brk(imm16(t, c.imm))
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
		back, derr := decodeOne(c.word)
		_ = derr
		require.Equal(t, in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}
}
