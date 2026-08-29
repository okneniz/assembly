package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestDbarCtor(t *testing.T) {
	c0, err := NewCode15(0)
	require.NoError(t, err)
	c1, err := NewCode15(1)
	require.NoError(t, err)

	// llvm-mc-verified: dbar 0 and dbar 1.
	require.Equal(t, uint32(0x38720000), ctorWord(t, NewDbar(c0)))
	require.Equal(t, uint32(0x38720001), ctorWord(t, NewDbar(c1)))
}

func TestDbarDecodeEncode(t *testing.T) {
	x, ok := decodeDbar(0x38720001, 0x90000000).(Dbar)
	require.True(t, ok, "type = %T, want Dbar", x)
	require.Equal(t, "dbar 1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(1), x.code.val)
	require.Equal(t, uint32(0x38720001), ctorWord(t, x))

	// llvm-mc-verified: dbar 0.
	n, ok2 := decodeDbar(0x38720000, 0x90000000).(Dbar)
	require.True(t, ok2, "type = %T, want Dbar", n)
	require.Equal(t, "dbar 0", n.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint32(0x38720000), ctorWord(t, n))
}
