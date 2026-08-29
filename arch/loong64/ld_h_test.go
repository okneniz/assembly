package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdHCtor(t *testing.T) {
	// llvm-mc-verified: ld.h $t0, $t1, 8.
	v, err := NewImm12(8)
	require.NoError(t, err)

	in := NewLdH(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x284021ac), ctorWord(t, in))

	_, ok := in.(LdH)
	require.True(t, ok, "type = %T, want LdH", in)
}

func TestLdHDecodeEncode(t *testing.T) {
	in := decodeLdH(0x284021ac, 0x90000000)

	x, ok := in.(LdH)
	require.True(t, ok, "type = %T, want LdH", in)
	require.Equal(t, "ld.h $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(8), x.imm.val)
	require.Equal(t, uint32(0x284021ac), ctorWord(t, x))

	// llvm-mc-verified: ld.h $t0, $t1, -16 (the negative byte offset
	// round-trips through the sign-extended field).
	y, ok := decodeLdH(0x287fc1ac, 0).(LdH)
	require.True(t, ok, "type = %T, want LdH", y)
	require.Equal(t, "ld.h $t0, $t1, -16", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-16), y.imm.val)
	require.Equal(t, uint32(0x287fc1ac), ctorWord(t, y))
}
