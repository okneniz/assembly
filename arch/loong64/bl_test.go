package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBlCtor(t *testing.T) {
	// llvm-mc-verified: bl 8 (at pc 0).
	in := NewBl(8)
	require.Equal(t, uint32(0x54000800), ctorWord(t, in))

	// llvm-mc-verified: bl -8 (at pc 0).
	neg := NewBl(-8)
	require.Equal(t, uint32(0x57fffbff), ctorWord(t, neg))
}

func TestBlDecodeEncode(t *testing.T) {
	x, ok := decodeBl(0x54000800, 0).(Bl)
	require.True(t, ok, "type = %T, want Bl", x)
	require.Equal(t, "bl 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(8), x.off.val)
	require.Equal(t, uint32(0x54000800), ctorWord(t, x))

	// The off field is the byte offset itself: the same word decodes
	// identically at any pc, and Encode is pc-independent.
	y, ok2 := decodeBl(0x54000800, 0x90000000).(Bl)
	require.True(t, ok2, "type = %T, want Bl", y)
	require.Equal(t, int64(8), y.off.val)
	require.Equal(t, uint32(0x54000800), ctorWord(t, y))

	// llvm-mc-verified: bl -8 - the split offs26 field round-trips (the
	// high d10 half lives in [9:0], the low k16 half in [25:10]).
	n, ok3 := decodeBl(0x57fffbff, 0).(Bl)
	require.True(t, ok3, "type = %T, want Bl", n)
	require.Equal(t, "bl -8", n.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), n.off.val)
	require.Equal(t, uint32(0x57fffbff), ctorWord(t, n))
}
