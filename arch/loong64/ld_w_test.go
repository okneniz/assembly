package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdWCtor(t *testing.T) {
	// llvm-mc-verified: ld.w $t0, $t1, 8.
	v, err := New().Imm12(8)
	require.NoError(t, err)

	in := New().LdW(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x288021ac), ctorWord(t, in))

	_, ok := in.(LdW)
	require.True(t, ok, "type = %T, want LdW", in)
}

func TestLdWDecodeEncode(t *testing.T) {
	in := decodeLdW(0x288021ac, 0x90000000)

	x, ok := in.(LdW)
	require.True(t, ok, "type = %T, want LdW", in)
	require.Equal(t, "ld.w $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(8), x.imm.val)
	require.Equal(t, uint32(0x288021ac), ctorWord(t, x))

	// llvm-mc-verified: ld.w $t0, $t1, -8 (the negative byte offset
	// round-trips through the sign-extended field).
	y, ok := decodeLdW(0x28bfe1ac, 0).(LdW)
	require.True(t, ok, "type = %T, want LdW", y)
	require.Equal(t, "ld.w $t0, $t1, -8", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), y.imm.val)
	require.Equal(t, uint32(0x28bfe1ac), ctorWord(t, y))
}
