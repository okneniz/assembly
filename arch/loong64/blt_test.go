package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBltCtor(t *testing.T) {
	// llvm-mc-verified: blt $t1, $t0, 8 (at pc 0; the manual order prints
	// rj first). The constructor takes (rj, rd, target).
	in := NewBlt(lreg(t, 13), lreg(t, 12), 8)
	require.Equal(t, uint32(0x600009ac), ctorWord(t, in))

	// llvm-mc-verified: blt $t1, $t0, -8 (at pc 0).
	neg := NewBlt(lreg(t, 13), lreg(t, 12), -8)
	require.Equal(t, uint32(0x63fff9ac), ctorWord(t, neg))
}

func TestBltDecodeEncode(t *testing.T) {
	x, ok := decodeBlt(0x600009ac, 0).(Blt)
	require.True(t, ok, "type = %T, want Blt", x)
	require.Equal(t, "blt $t1, $t0, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(8), x.off.val)
	require.Equal(t, uint32(0x600009ac), ctorWord(t, x))

	// The off field is the byte offset itself: the same word decodes
	// identically at any pc, and Encode is pc-independent.
	y, ok2 := decodeBlt(0x600009ac, 0x90000000).(Blt)
	require.True(t, ok2, "type = %T, want Blt", y)
	require.Equal(t, int64(8), y.off.val)
	require.Equal(t, uint32(0x600009ac), ctorWord(t, y))

	// llvm-mc-verified: blt $t1, $t0, -8 - the negative target round-trips.
	n, ok3 := decodeBlt(0x63fff9ac, 0).(Blt)
	require.True(t, ok3, "type = %T, want Blt", n)
	require.Equal(t, "blt $t1, $t0, -8", n.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), n.off.val)
	require.Equal(t, uint32(0x63fff9ac), ctorWord(t, n))
}
