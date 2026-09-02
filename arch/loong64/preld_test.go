package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestPreldCtor(t *testing.T) {
	// llvm-mc-verified: preld 5, $t1, 8 (the manual prints the hint
	// first).
	h, err := New().UImm5(5)
	require.NoError(t, err)
	off, err := New().Imm12(8)
	require.NoError(t, err)

	in := New().Preld(h, lreg(t, 13), off)
	require.Equal(t, uint32(0x2ac021a5), ctorWord(t, in))

	_, ok := in.(Preld)
	require.True(t, ok, "type = %T, want Preld", in)
}

func TestPreldDecodeEncode(t *testing.T) {
	in := decodePreld(0x2ac021a5, 0x90000000)

	x, ok := in.(Preld)
	require.True(t, ok, "type = %T, want Preld", in)
	require.Equal(t, "preld 5, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(5), x.hint.val)
	require.Equal(t, int64(8), x.off.val)
	require.Equal(t, uint32(0x2ac021a5), ctorWord(t, x))

	// llvm-mc-verified: preld 0, $t1, -8 (the negative byte offset
	// round-trips through the sign-extended field).
	y, ok := decodePreld(0x2affe1a0, 0).(Preld)
	require.True(t, ok, "type = %T, want Preld", y)
	require.Equal(t, "preld 0, $t1, -8", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint32(0x2affe1a0), ctorWord(t, y))
}
