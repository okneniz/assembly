package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStWCtor(t *testing.T) {
	// llvm-mc-verified: st.w $t0, $t1, 8.
	v, err := New().Imm12(8)
	require.NoError(t, err)

	in := New().StW(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x298021ac), ctorWord(t, in))

	_, ok := in.(StW)
	require.True(t, ok, "type = %T, want StW", in)
}

func TestStWDecodeEncode(t *testing.T) {
	in := decodeStW(0x298021ac, 0x90000000)

	x, ok := in.(StW)
	require.True(t, ok, "type = %T, want StW", in)
	require.Equal(t, "st.w $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(8), x.imm.val)
	require.Equal(t, uint32(0x298021ac), ctorWord(t, x))

	// llvm-mc-verified: st.w $t0, $t1, -8 (the negative byte offset
	// round-trips through the sign-extended field).
	y, ok := decodeStW(0x29bfe1ac, 0).(StW)
	require.True(t, ok, "type = %T, want StW", y)
	require.Equal(t, "st.w $t0, $t1, -8", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), y.imm.val)
	require.Equal(t, uint32(0x29bfe1ac), ctorWord(t, y))
}
