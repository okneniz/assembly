package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBltuCtor(t *testing.T) {
	// llvm-mc-verified: bltu $t1, $t0, 8 (at pc 0; the manual order prints
	// rj first). The constructor takes (rj, rd, target).
	in := NewBltu(lreg(t, 13), lreg(t, 12), 8)
	require.Equal(t, uint32(0x680009ac), ctorWord(t, in))

	// llvm-mc-verified: bltu $t1, $t0, -8 (at pc 0).
	neg := NewBltu(lreg(t, 13), lreg(t, 12), -8)
	require.Equal(t, uint32(0x6bfff9ac), ctorWord(t, neg))
}

func TestBltuDecodeEncode(t *testing.T) {
	x, ok := decodeBltu(0x680009ac, 0).(Bltu)
	require.True(t, ok, "type = %T, want Bltu", x)
	require.Equal(t, "bltu $t1, $t0, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(8), x.off.val)
	require.Equal(t, uint32(0x680009ac), ctorWord(t, x))

	// The off field is the byte offset itself: the same word decodes
	// identically at any pc, and Encode is pc-independent.
	y, ok2 := decodeBltu(0x680009ac, 0x90000000).(Bltu)
	require.True(t, ok2, "type = %T, want Bltu", y)
	require.Equal(t, int64(8), y.off.val)
	require.Equal(t, uint32(0x680009ac), ctorWord(t, y))

	// llvm-mc-verified: bltu $t1, $t0, -8 - the negative target round-trips.
	n, ok3 := decodeBltu(0x6bfff9ac, 0).(Bltu)
	require.True(t, ok3, "type = %T, want Bltu", n)
	require.Equal(t, "bltu $t1, $t0, -8", n.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), n.off.val)
	require.Equal(t, uint32(0x6bfff9ac), ctorWord(t, n))
}
