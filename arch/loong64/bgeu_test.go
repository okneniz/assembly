package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBgeuCtor(t *testing.T) {
	// llvm-mc-verified: bgeu $t1, $t0, 8 (at pc 0; the manual order prints
	// rj first). The constructor takes (rj, rd, target).
	in := New().Bgeu(lreg(t, 13), lreg(t, 12), 8)
	require.Equal(t, uint32(0x6c0009ac), ctorWord(t, in))

	// llvm-mc-verified: bgeu $t1, $t0, -8 (at pc 0).
	neg := New().Bgeu(lreg(t, 13), lreg(t, 12), -8)
	require.Equal(t, uint32(0x6ffff9ac), ctorWord(t, neg))
}

func TestBgeuDecodeEncode(t *testing.T) {
	x, ok := decodeBgeu(0x6c0009ac, 0).(Bgeu)
	require.True(t, ok, "type = %T, want Bgeu", x)
	require.Equal(t, "bgeu $t1, $t0, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(8), x.off.val)
	require.Equal(t, uint32(0x6c0009ac), ctorWord(t, x))

	// The off field is the byte offset itself: the same word decodes
	// identically at any pc, and Encode is pc-independent.
	y, ok2 := decodeBgeu(0x6c0009ac, 0x90000000).(Bgeu)
	require.True(t, ok2, "type = %T, want Bgeu", y)
	require.Equal(t, int64(8), y.off.val)
	require.Equal(t, uint32(0x6c0009ac), ctorWord(t, y))

	// llvm-mc-verified: bgeu $t1, $t0, -8 - the negative target round-trips.
	n, ok3 := decodeBgeu(0x6ffff9ac, 0).(Bgeu)
	require.True(t, ok3, "type = %T, want Bgeu", n)
	require.Equal(t, "bgeu $t1, $t0, -8", n.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), n.off.val)
	require.Equal(t, uint32(0x6ffff9ac), ctorWord(t, n))
}
