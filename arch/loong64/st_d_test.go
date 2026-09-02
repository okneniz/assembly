package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStDCtor(t *testing.T) {
	// llvm-mc-verified: st.d $t0, $t1, 8.
	v, err := New().Imm12(8)
	require.NoError(t, err)

	in := New().StD(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x29c021ac), ctorWord(t, in))

	_, ok := in.(StD)
	require.True(t, ok, "type = %T, want StD", in)
}

func TestStDDecodeEncode(t *testing.T) {
	in := decodeStD(0x29c021ac, 0x90000000)

	x, ok := in.(StD)
	require.True(t, ok, "type = %T, want StD", in)
	require.Equal(t, "st.d $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(8), x.imm.val)
	require.Equal(t, uint32(0x29c021ac), ctorWord(t, x))

	// llvm-mc-verified: st.d $t0, $t1, -8 (the negative byte offset
	// round-trips through the sign-extended field).
	y, ok := decodeStD(0x29ffe1ac, 0).(StD)
	require.True(t, ok, "type = %T, want StD", y)
	require.Equal(t, "st.d $t0, $t1, -8", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), y.imm.val)
	require.Equal(t, uint32(0x29ffe1ac), ctorWord(t, y))
}
