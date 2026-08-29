package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStBCtor(t *testing.T) {
	// llvm-mc-verified: st.b $t0, $t1, 8.
	v, err := NewImm12(8)
	require.NoError(t, err)

	in := NewStB(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x290021ac), ctorWord(t, in))

	_, ok := in.(StB)
	require.True(t, ok, "type = %T, want StB", in)
}

func TestStBDecodeEncode(t *testing.T) {
	in := decodeStB(0x290021ac, 0x90000000)

	x, ok := in.(StB)
	require.True(t, ok, "type = %T, want StB", in)
	require.Equal(t, "st.b $t0, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(8), x.imm.val)
	require.Equal(t, uint32(0x290021ac), ctorWord(t, x))

	// llvm-mc-verified: st.b $t0, $t1, -8 (the negative byte offset
	// round-trips through the sign-extended field).
	y, ok := decodeStB(0x293fe1ac, 0).(StB)
	require.True(t, ok, "type = %T, want StB", y)
	require.Equal(t, "st.b $t0, $t1, -8", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, int64(-8), y.imm.val)
	require.Equal(t, uint32(0x293fe1ac), ctorWord(t, y))
}
