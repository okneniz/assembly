package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLlacqWCtor(t *testing.T) {
	// llvm-mc-verified: llacq.w $t0, $t1.
	require.Equal(
		t,
		uint32(0x385781ac),
		ctorWord(t, NewLlacqW(lreg(t, 12), lreg(t, 13))),
	)

	in := NewLlacqW(lreg(t, 1), lreg(t, 2))
	_, ok := in.(LlacqW)
	require.True(t, ok, "type = %T, want LlacqW", in)
}

func TestLlacqWDecodeEncode(t *testing.T) {
	in := decodeLlacqW(0x385781ac, 0x90000000)

	llacqw, ok := in.(LlacqW)
	require.True(t, ok, "type = %T, want LlacqW", in)
	require.Equal(t, "llacq.w $t0, $t1", llacqw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), llacqw.Addr())
	require.Equal(t, 4, llacqw.Len())
	require.Equal(t, uint32(0x385781ac), ctorWord(t, llacqw))
}
