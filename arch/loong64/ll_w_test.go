package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLlWCtor(t *testing.T) {
	// llvm-mc-verified: ll.w $t0, $t1, 8.
	off, err := NewImm14(8)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x200009ac),
		ctorWord(t, NewLlW(lreg(t, 12), lreg(t, 13), off)),
	)

	in := NewLlW(lreg(t, 1), lreg(t, 2), off)
	_, ok := in.(LlW)
	require.True(t, ok, "type = %T, want LlW", in)
}

func TestLlWDecodeEncode(t *testing.T) {
	in := decodeLlW(0x200009ac, 0x90000000)

	llw, ok := in.(LlW)
	require.True(t, ok, "type = %T, want LlW", in)
	require.Equal(t, "ll.w $t0, $t1, 8", llw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), llw.Addr())
	require.Equal(t, 4, llw.Len())
	require.Equal(t, int64(8), llw.off.val)
	require.Equal(t, uint32(0x200009ac), ctorWord(t, llw))
}
