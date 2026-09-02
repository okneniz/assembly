package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBreakCtor(t *testing.T) {
	c0, err := New().Code15(0)
	require.NoError(t, err)
	c1, err := New().Code15(1)
	require.NoError(t, err)

	// llvm-mc-verified: break 0 and break 1.
	require.Equal(t, uint32(0x002a0000), ctorWord(t, New().Break(c0)))
	require.Equal(t, uint32(0x002a0001), ctorWord(t, New().Break(c1)))
}

func TestBreakDecodeEncode(t *testing.T) {
	x, ok := decodeBreak(0x002a0001, 0x90000000).(Break)
	require.True(t, ok, "type = %T, want Break", x)
	require.Equal(t, "break 1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(1), x.code.val)
	require.Equal(t, uint32(0x002a0001), ctorWord(t, x))

	// llvm-mc-verified: break 0.
	n, ok2 := decodeBreak(0x002a0000, 0x90000000).(Break)
	require.True(t, ok2, "type = %T, want Break", n)
	require.Equal(t, "break 0", n.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint32(0x002a0000), ctorWord(t, n))
}
