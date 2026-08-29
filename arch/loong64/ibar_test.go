package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestIbarCtor(t *testing.T) {
	c0, err := NewCode15(0)
	require.NoError(t, err)

	// llvm-mc-verified: ibar 0.
	require.Equal(t, uint32(0x38728000), ctorWord(t, NewIbar(c0)))
}

func TestIbarDecodeEncode(t *testing.T) {
	x, ok := decodeIbar(0x38728000, 0x90000000).(Ibar)
	require.True(t, ok, "type = %T, want Ibar", x)
	require.Equal(t, "ibar 0", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(0), x.code.val)
	require.Equal(t, uint32(0x38728000), ctorWord(t, x))
}
