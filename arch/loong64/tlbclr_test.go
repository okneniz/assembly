package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestTlbclrCtor(t *testing.T) {
	// llvm-mc-verified: tlbclr.
	require.Equal(t, uint32(0x06482000), ctorWord(t, New().Tlbclr()))

	in := New().Tlbclr()
	_, ok := in.(Tlbclr)
	require.True(t, ok, "type = %T, want Tlbclr", in)
}

func TestTlbclrDecodeEncode(t *testing.T) {
	// llvm-mc-verified: tlbclr.
	in := decodeTlbclr(0x06482000, 0x90000000)

	x, ok := in.(Tlbclr)
	require.True(t, ok, "type = %T, want Tlbclr", in)
	require.Equal(t, "tlbclr", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x06482000), ctorWord(t, x))
}
