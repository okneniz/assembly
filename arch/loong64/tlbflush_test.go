package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestTlbflushCtor(t *testing.T) {
	// llvm-mc-verified: tlbflush.
	require.Equal(t, uint32(0x06482400), ctorWord(t, New().Tlbflush()))

	in := New().Tlbflush()
	_, ok := in.(Tlbflush)
	require.True(t, ok, "type = %T, want Tlbflush", in)
}

func TestTlbflushDecodeEncode(t *testing.T) {
	// llvm-mc-verified: tlbflush.
	in := decodeTlbflush(0x06482400, 0x90000000)

	x, ok := in.(Tlbflush)
	require.True(t, ok, "type = %T, want Tlbflush", in)
	require.Equal(t, "tlbflush", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x06482400), ctorWord(t, x))
}
