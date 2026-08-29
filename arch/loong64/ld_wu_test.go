package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdWuCtor(t *testing.T) {
	// llvm-mc-verified: ld.wu $t0, $t1, 8.
	v, err := NewImm12(8)
	require.NoError(t, err)

	in := NewLdWu(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x2a8021ac), ctorWord(t, in))

	_, ok := in.(LdWu)
	require.True(t, ok, "type = %T, want LdWu", in)
}

func TestLdWuDecodeEncode(t *testing.T) {
	in := decodeLdWu(0x2a8021ac, 0x90000000)

	x, ok := in.(LdWu)
	require.True(t, ok, "type = %T, want LdWu", in)
	require.Equal(t, "ld.wu $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(8), x.imm.val)
	require.Equal(t, uint32(0x2a8021ac), ctorWord(t, x))

	// llvm-mc-verified: ld.wu $t0, $t1, -8 (the negative byte offset
	// round-trips through the sign-extended field).
	y, ok := decodeLdWu(0x2abfe1ac, 0).(LdWu)
	require.True(t, ok, "type = %T, want LdWu", y)
	require.Equal(t, "ld.wu $t0, $t1, -8", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), y.imm.val)
	require.Equal(t, uint32(0x2abfe1ac), ctorWord(t, y))
}
