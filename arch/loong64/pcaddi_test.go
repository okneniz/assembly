package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestPcaddiCtor(t *testing.T) {
	// llvm-mc-verified: pcaddi $t0, 5.
	v, err := New().Imm20(5)
	require.NoError(t, err)

	in := New().Pcaddi(lreg(t, 12), v)
	require.Equal(t, uint32(0x180000ac), ctorWord(t, in))

	_, ok := in.(Pcaddi)
	require.True(t, ok, "type = %T, want Pcaddi", in)
}

func TestPcaddiDecodeEncode(t *testing.T) {
	in := decodePcaddi(0x180000ac, 0x90000000)

	x, ok := in.(Pcaddi)
	require.True(t, ok, "type = %T, want Pcaddi", in)
	require.Equal(t, "pcaddi $t0, 5", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(5), x.imm.val)
	require.Equal(t, uint32(0x180000ac), ctorWord(t, x))

	// llvm-mc-verified: pcaddi $t0, -1 (the raw si20 round-trips).
	y, ok := decodePcaddi(0x19ffffec, 0).(Pcaddi)
	require.True(t, ok, "type = %T, want Pcaddi", y)
	require.Equal(t, "pcaddi $t0, -1", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint32(0x19ffffec), ctorWord(t, y))
}
