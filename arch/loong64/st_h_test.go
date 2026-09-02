package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStHCtor(t *testing.T) {
	// llvm-mc-verified: st.h $t0, $t1, 8.
	v, err := New().Imm12(8)
	require.NoError(t, err)

	in := New().StH(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x294021ac), ctorWord(t, in))

	_, ok := in.(StH)
	require.True(t, ok, "type = %T, want StH", in)
}

func TestStHDecodeEncode(t *testing.T) {
	in := decodeStH(0x294021ac, 0x90000000)

	x, ok := in.(StH)
	require.True(t, ok, "type = %T, want StH", in)
	require.Equal(t, "st.h $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(8), x.imm.val)
	require.Equal(t, uint32(0x294021ac), ctorWord(t, x))

	// llvm-mc-verified: st.h $t0, $t1, -8 (the negative byte offset
	// round-trips through the sign-extended field).
	y, ok := decodeStH(0x297fe1ac, 0).(StH)
	require.True(t, ok, "type = %T, want StH", y)
	require.Equal(t, "st.h $t0, $t1, -8", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), y.imm.val)
	require.Equal(t, uint32(0x297fe1ac), ctorWord(t, y))
}
