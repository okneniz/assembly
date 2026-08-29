package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLu32iDCtor(t *testing.T) {
	// llvm-mc-verified: lu32i.d $t0, 5.
	v, err := NewImm20(5)
	require.NoError(t, err)

	in := NewLu32iD(lreg(t, 12), v)
	require.Equal(t, uint32(0x160000ac), ctorWord(t, in))

	_, ok := in.(Lu32iD)
	require.True(t, ok, "type = %T, want Lu32iD", in)
}

func TestLu32iDDecodeEncode(t *testing.T) {
	in := decodeLu32iD(0x160000ac, 0x90000000)

	x, ok := in.(Lu32iD)
	require.True(t, ok, "type = %T, want Lu32iD", in)
	require.Equal(t, "lu32i.d $t0, 5", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(5), x.imm.val)
	require.Equal(t, uint32(0x160000ac), ctorWord(t, x))

	// llvm-mc-verified: lu32i.d $t0, -1 (the raw si20 round-trips).
	y, ok := decodeLu32iD(0x17ffffec, 0).(Lu32iD)
	require.True(t, ok, "type = %T, want Lu32iD", y)
	require.Equal(t, "lu32i.d $t0, -1", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint32(0x17ffffec), ctorWord(t, y))
}
