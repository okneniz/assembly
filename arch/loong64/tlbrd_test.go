package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestTlbrdCtor(t *testing.T) {
	// llvm-mc-verified: tlbrd.
	require.Equal(t, uint32(0x06482c00), ctorWord(t, New().Tlbrd()))

	in := New().Tlbrd()
	_, ok := in.(Tlbrd)
	require.True(t, ok, "type = %T, want Tlbrd", in)
}

func TestTlbrdDecodeEncode(t *testing.T) {
	// llvm-mc-verified: tlbrd.
	in := decodeTlbrd(0x06482c00, 0x90000000)

	x, ok := in.(Tlbrd)
	require.True(t, ok, "type = %T, want Tlbrd", in)
	require.Equal(t, "tlbrd", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x06482c00), ctorWord(t, x))
}
