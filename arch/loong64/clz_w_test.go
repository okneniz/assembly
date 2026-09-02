package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestClzWCtor(t *testing.T) {
	// llvm-mc-verified: clz.w $t0, $t1.
	require.Equal(
		t,
		uint32(0x000015ac),
		ctorWord(t, New().ClzW(lreg(t, 12), lreg(t, 13))),
	)

	in := New().ClzW(lreg(t, 1), lreg(t, 2))
	_, ok := in.(ClzW)
	require.True(t, ok, "type = %T, want ClzW", in)
}

func TestClzWDecodeEncode(t *testing.T) {
	in := decodeClzW(0x000015ac, 0x90000000)

	clzw, ok := in.(ClzW)
	require.True(t, ok, "type = %T, want ClzW", in)
	require.Equal(t, "clz.w $t0, $t1", clzw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), clzw.Addr())
	require.Equal(t, uint32(0x000015ac), ctorWord(t, clzw))
}
