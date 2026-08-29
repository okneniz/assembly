package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBneCtor(t *testing.T) {
	// llvm-mc-verified: bne $t1, $t0, 8 (at pc 0; the manual order prints
	// rj first). The constructor takes (rj, rd, target).
	in := NewBne(lreg(t, 13), lreg(t, 12), 8)
	require.Equal(t, uint32(0x5c0009ac), ctorWord(t, in))

	// llvm-mc-verified: bne $t1, $t0, -8 (at pc 0).
	neg := NewBne(lreg(t, 13), lreg(t, 12), -8)
	require.Equal(t, uint32(0x5ffff9ac), ctorWord(t, neg))
}

func TestBneDecodeEncode(t *testing.T) {
	x, ok := decodeBne(0x5c0009ac, 0).(Bne)
	require.True(t, ok, "type = %T, want Bne", x)
	require.Equal(t, "bne $t1, $t0, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(8), x.off.val)
	require.Equal(t, uint32(0x5c0009ac), ctorWord(t, x))

	// The off field is the byte offset itself: the same word decodes
	// identically at any pc, and Encode is pc-independent.
	y, ok2 := decodeBne(0x5c0009ac, 0x90000000).(Bne)
	require.True(t, ok2, "type = %T, want Bne", y)
	require.Equal(t, int64(8), y.off.val)
	require.Equal(t, uint32(0x5c0009ac), ctorWord(t, y))

	// llvm-mc-verified: bne $t1, $t0, -8 - the negative target round-trips.
	n, ok3 := decodeBne(0x5ffff9ac, 0).(Bne)
	require.True(t, ok3, "type = %T, want Bne", n)
	require.Equal(t, "bne $t1, $t0, -8", n.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), n.off.val)
	require.Equal(t, uint32(0x5ffff9ac), ctorWord(t, n))
}
