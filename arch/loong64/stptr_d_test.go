package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStptrDCtor(t *testing.T) {
	// llvm-mc-verified: stptr.d $t0, $t1, 8 (the si14 word count is 2).
	v, err := New().Imm14(8)
	require.NoError(t, err)

	in := New().StptrD(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x270009ac), ctorWord(t, in))

	_, ok := in.(StptrD)
	require.True(t, ok, "type = %T, want StptrD", in)
}

func TestStptrDDecodeEncode(t *testing.T) {
	in := decodeStptrD(0x270009ac, 0x90000000)

	x, ok := in.(StptrD)
	require.True(t, ok, "type = %T, want StptrD", in)
	require.Equal(t, "stptr.d $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(8), x.off.val)
	require.Equal(t, uint32(0x270009ac), ctorWord(t, x))

	// llvm-mc-verified: stptr.d $t0, $t1, -8 (the byte offset is stored
	// raw, the si14 word count -2 round-trips).
	y, ok := decodeStptrD(0x27fff9ac, 0).(StptrD)
	require.True(t, ok, "type = %T, want StptrD", y)
	require.Equal(t, "stptr.d $t0, $t1, -8", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), y.off.val)
	require.Equal(t, uint32(0x27fff9ac), ctorWord(t, y))
}
