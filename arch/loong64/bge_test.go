package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBgeCtor(t *testing.T) {
	// llvm-mc-verified: bge $t1, $t0, 8 (at pc 0; the manual order prints
	// rj first). The constructor takes (rj, rd, target).
	in := New().Bge(lreg(t, 13), lreg(t, 12), 8)
	require.Equal(t, uint32(0x640009ac), ctorWord(t, in))

	// llvm-mc-verified: bge $t1, $t0, -8 (at pc 0).
	neg := New().Bge(lreg(t, 13), lreg(t, 12), -8)
	require.Equal(t, uint32(0x67fff9ac), ctorWord(t, neg))
}

func TestBgeDecodeEncode(t *testing.T) {
	x, ok := decodeBge(0x640009ac, 0).(Bge)
	require.True(t, ok, "type = %T, want Bge", x)
	require.Equal(t, "bge $t1, $t0, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(8), x.off.val)
	require.Equal(t, uint32(0x640009ac), ctorWord(t, x))

	// The off field is the byte offset itself: the same word decodes
	// identically at any pc, and Encode is pc-independent.
	y, ok2 := decodeBge(0x640009ac, 0x90000000).(Bge)
	require.True(t, ok2, "type = %T, want Bge", y)
	require.Equal(t, int64(8), y.off.val)
	require.Equal(t, uint32(0x640009ac), ctorWord(t, y))

	// llvm-mc-verified: bge $t1, $t0, -8 - the negative target round-trips.
	n, ok3 := decodeBge(0x67fff9ac, 0).(Bge)
	require.True(t, ok3, "type = %T, want Bge", n)
	require.Equal(t, "bge $t1, $t0, -8", n.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), n.off.val)
	require.Equal(t, uint32(0x67fff9ac), ctorWord(t, n))
}
