package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestTlbwrCtor(t *testing.T) {
	// llvm-mc-verified: tlbwr.
	require.Equal(t, uint32(0x06483000), ctorWord(t, NewTlbwr()))

	in := NewTlbwr()
	_, ok := in.(Tlbwr)
	require.True(t, ok, "type = %T, want Tlbwr", in)
}

func TestTlbwrDecodeEncode(t *testing.T) {
	// llvm-mc-verified: tlbwr.
	in := decodeTlbwr(0x06483000, 0x90000000)

	x, ok := in.(Tlbwr)
	require.True(t, ok, "type = %T, want Tlbwr", in)
	require.Equal(t, "tlbwr", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x06483000), ctorWord(t, x))
}
