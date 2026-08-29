package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdBCtor(t *testing.T) {
	// llvm-mc-verified: ld.b $t0, $t1, 8.
	v, err := NewImm12(8)
	require.NoError(t, err)

	in := NewLdB(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x280021ac), ctorWord(t, in))

	_, ok := in.(LdB)
	require.True(t, ok, "type = %T, want LdB", in)
}

func TestLdBDecodeEncode(t *testing.T) {
	in := decodeLdB(0x280021ac, 0x90000000)

	x, ok := in.(LdB)
	require.True(t, ok, "type = %T, want LdB", in)
	require.Equal(t, "ld.b $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(8), x.imm.val)
	require.Equal(t, uint32(0x280021ac), ctorWord(t, x))

	// llvm-mc-verified: ld.b $t0, $t1, -8 (the negative byte offset
	// round-trips through the sign-extended field).
	y, ok := decodeLdB(0x283fe1ac, 0).(LdB)
	require.True(t, ok, "type = %T, want LdB", y)
	require.Equal(t, "ld.b $t0, $t1, -8", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), y.imm.val)
	require.Equal(t, uint32(0x283fe1ac), ctorWord(t, y))
}
