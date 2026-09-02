package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLlacqDCtor(t *testing.T) {
	// llvm-mc-verified: llacq.d $t0, $t1.
	require.Equal(
		t,
		uint32(0x385789ac),
		ctorWord(t, New().LlacqD(lreg(t, 12), lreg(t, 13))),
	)

	in := New().LlacqD(lreg(t, 1), lreg(t, 2))
	_, ok := in.(LlacqD)
	require.True(t, ok, "type = %T, want LlacqD", in)
}

func TestLlacqDDecodeEncode(t *testing.T) {
	in := decodeLlacqD(0x385789ac, 0x90000000)

	llacqd, ok := in.(LlacqD)
	require.True(t, ok, "type = %T, want LlacqD", in)
	require.Equal(t, "llacq.d $t0, $t1", llacqd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), llacqd.Addr())
	require.Equal(t, 4, llacqd.Len())
	require.Equal(t, uint32(0x385789ac), ctorWord(t, llacqd))
}
