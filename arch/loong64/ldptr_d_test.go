package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdptrDCtor(t *testing.T) {
	// llvm-mc-verified: ldptr.d $t0, $t1, 8 (the si14 word count is 2).
	v, err := NewImm14(8)
	require.NoError(t, err)

	in := NewLdptrD(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x260009ac), ctorWord(t, in))

	_, ok := in.(LdptrD)
	require.True(t, ok, "type = %T, want LdptrD", in)
}

func TestLdptrDDecodeEncode(t *testing.T) {
	in := decodeLdptrD(0x260009ac, 0x90000000)

	x, ok := in.(LdptrD)
	require.True(t, ok, "type = %T, want LdptrD", in)
	require.Equal(t, "ldptr.d $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(8), x.off.val)
	require.Equal(t, uint32(0x260009ac), ctorWord(t, x))

	// llvm-mc-verified: ldptr.d $t0, $t1, -8 (the byte offset is stored
	// raw, the si14 word count -2 round-trips).
	y, ok := decodeLdptrD(0x26fff9ac, 0).(LdptrD)
	require.True(t, ok, "type = %T, want LdptrD", y)
	require.Equal(t, "ldptr.d $t0, $t1, -8", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), y.off.val)
	require.Equal(t, uint32(0x26fff9ac), ctorWord(t, y))
}
