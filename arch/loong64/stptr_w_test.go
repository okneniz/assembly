package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStptrWCtor(t *testing.T) {
	// llvm-mc-verified: stptr.w $t0, $t1, 8 (the si14 word count is 2).
	v, err := New().Imm14(8)
	require.NoError(t, err)

	in := New().StptrW(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x250009ac), ctorWord(t, in))

	_, ok := in.(StptrW)
	require.True(t, ok, "type = %T, want StptrW", in)
}

func TestStptrWDecodeEncode(t *testing.T) {
	in := decodeStptrW(0x250009ac, 0x90000000)

	x, ok := in.(StptrW)
	require.True(t, ok, "type = %T, want StptrW", in)
	require.Equal(t, "stptr.w $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(8), x.off.val)
	require.Equal(t, uint32(0x250009ac), ctorWord(t, x))

	// llvm-mc-verified: stptr.w $t0, $t1, -8 (the byte offset is stored
	// raw, the si14 word count -2 round-trips).
	y, ok := decodeStptrW(0x25fff9ac, 0).(StptrW)
	require.True(t, ok, "type = %T, want StptrW", y)
	require.Equal(t, "stptr.w $t0, $t1, -8", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), y.off.val)
	require.Equal(t, uint32(0x25fff9ac), ctorWord(t, y))
}
