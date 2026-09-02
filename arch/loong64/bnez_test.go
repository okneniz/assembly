package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBnezCtor(t *testing.T) {
	// llvm-mc-verified: bnez $t1, 8 (at pc 0).
	in := New().Bnez(lreg(t, 13), 8)
	require.Equal(t, uint32(0x440009a0), ctorWord(t, in))

	// llvm-mc-verified: bnez $t1, -8 (at pc 0).
	neg := New().Bnez(lreg(t, 13), -8)
	require.Equal(t, uint32(0x47fff9bf), ctorWord(t, neg))
}

func TestBnezDecodeEncode(t *testing.T) {
	x, ok := decodeBnez(0x440009a0, 0).(Bnez)
	require.True(t, ok, "type = %T, want Bnez", x)
	require.Equal(t, "bnez $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(8), x.off.val)
	require.Equal(t, uint32(0x440009a0), ctorWord(t, x))

	// The split offs21 field: a negative offset round-trips (the high d5
	// half lives in [4:0], the low k16 half in [25:10]); the off field
	// is the offset itself, identical at any pc.
	n, ok2 := decodeBnez(0x47fff9bf, 0x1000).(Bnez)
	require.True(t, ok2, "type = %T, want Bnez", n)
	require.Equal(t, "bnez $t1, -8", n.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), n.off.val)
	require.Equal(t, uint32(0x47fff9bf), ctorWord(t, n))
}
