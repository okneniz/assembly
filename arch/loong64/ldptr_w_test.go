package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdptrWCtor(t *testing.T) {
	// llvm-mc-verified: ldptr.w $t0, $t1, 8 (the si14 word count is 2).
	v, err := NewImm14(8)
	require.NoError(t, err)

	in := NewLdptrW(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x240009ac), ctorWord(t, in))

	_, ok := in.(LdptrW)
	require.True(t, ok, "type = %T, want LdptrW", in)
}

func TestLdptrWDecodeEncode(t *testing.T) {
	in := decodeLdptrW(0x240009ac, 0x90000000)

	x, ok := in.(LdptrW)
	require.True(t, ok, "type = %T, want LdptrW", in)
	require.Equal(t, "ldptr.w $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(8), x.off.val)
	require.Equal(t, uint32(0x240009ac), ctorWord(t, x))

	// llvm-mc-verified: ldptr.w $t0, $t1, -8 (the byte offset is stored
	// raw, the si14 word count -2 round-trips).
	y, ok := decodeLdptrW(0x24fff9ac, 0).(LdptrW)
	require.True(t, ok, "type = %T, want LdptrW", y)
	require.Equal(t, "ldptr.w $t0, $t1, -8", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), y.off.val)
	require.Equal(t, uint32(0x24fff9ac), ctorWord(t, y))
}
