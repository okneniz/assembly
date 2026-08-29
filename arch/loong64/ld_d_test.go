package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdDCtor(t *testing.T) {
	// llvm-mc-verified: ld.d $t0, $t1, 8.
	v, err := NewImm12(8)
	require.NoError(t, err)

	in := NewLdD(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x28c021ac), ctorWord(t, in))

	_, ok := in.(LdD)
	require.True(t, ok, "type = %T, want LdD", in)
}

func TestLdDDecodeEncode(t *testing.T) {
	in := decodeLdD(0x28c021ac, 0x90000000)

	x, ok := in.(LdD)
	require.True(t, ok, "type = %T, want LdD", in)
	require.Equal(t, "ld.d $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(8), x.imm.val)
	require.Equal(t, uint32(0x28c021ac), ctorWord(t, x))

	// llvm-mc-verified: ld.d $t0, $t1, -8 (the negative byte offset
	// round-trips through the sign-extended field).
	y, ok := decodeLdD(0x28ffe1ac, 0).(LdD)
	require.True(t, ok, "type = %T, want LdD", y)
	require.Equal(t, "ld.d $t0, $t1, -8", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), y.imm.val)
	require.Equal(t, uint32(0x28ffe1ac), ctorWord(t, y))
}
