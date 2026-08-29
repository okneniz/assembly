package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestTlbsrchCtor(t *testing.T) {
	// llvm-mc-verified: tlbsrch.
	require.Equal(t, uint32(0x06482800), ctorWord(t, NewTlbsrch()))

	in := NewTlbsrch()
	_, ok := in.(Tlbsrch)
	require.True(t, ok, "type = %T, want Tlbsrch", in)
}

func TestTlbsrchDecodeEncode(t *testing.T) {
	// llvm-mc-verified: tlbsrch.
	in := decodeTlbsrch(0x06482800, 0x90000000)

	x, ok := in.(Tlbsrch)
	require.True(t, ok, "type = %T, want Tlbsrch", in)
	require.Equal(t, "tlbsrch", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x06482800), ctorWord(t, x))
}
