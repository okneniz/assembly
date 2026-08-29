package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLlDCtor(t *testing.T) {
	// llvm-mc-verified: ll.d $t0, $t1, 8.
	off, err := NewImm14(8)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x220009ac),
		ctorWord(t, NewLlD(lreg(t, 12), lreg(t, 13), off)),
	)

	in := NewLlD(lreg(t, 1), lreg(t, 2), off)
	_, ok := in.(LlD)
	require.True(t, ok, "type = %T, want LlD", in)
}

func TestLlDDecodeEncode(t *testing.T) {
	in := decodeLlD(0x220009ac, 0x90000000)

	lld, ok := in.(LlD)
	require.True(t, ok, "type = %T, want LlD", in)
	require.Equal(t, "ll.d $t0, $t1, 8", lld.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), lld.Addr())
	require.Equal(t, 4, lld.Len())
	require.Equal(t, int64(8), lld.off.val)
	require.Equal(t, uint32(0x220009ac), ctorWord(t, lld))
}
