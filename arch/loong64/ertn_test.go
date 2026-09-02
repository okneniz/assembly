package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestErtnCtor(t *testing.T) {
	// llvm-mc-verified: ertn.
	require.Equal(t, uint32(0x06483800), ctorWord(t, New().Ertn()))

	in := New().Ertn()
	_, ok := in.(Ertn)
	require.True(t, ok, "type = %T, want Ertn", in)
}

func TestErtnDecodeEncode(t *testing.T) {
	// llvm-mc-verified: ertn.
	in := decodeErtn(0x06483800, 0x90000000)

	x, ok := in.(Ertn)
	require.True(t, ok, "type = %T, want Ertn", in)
	require.Equal(t, "ertn", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x06483800), ctorWord(t, x))
}
