package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestTlbfillCtor(t *testing.T) {
	// llvm-mc-verified: tlbfill.
	require.Equal(t, uint32(0x06483400), ctorWord(t, New().Tlbfill()))

	in := New().Tlbfill()
	_, ok := in.(Tlbfill)
	require.True(t, ok, "type = %T, want Tlbfill", in)
}

func TestTlbfillDecodeEncode(t *testing.T) {
	// llvm-mc-verified: tlbfill.
	in := decodeTlbfill(0x06483400, 0x90000000)

	x, ok := in.(Tlbfill)
	require.True(t, ok, "type = %T, want Tlbfill", in)
	require.Equal(t, "tlbfill", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x06483400), ctorWord(t, x))
}
