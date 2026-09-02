package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestScDCtor(t *testing.T) {
	// llvm-mc-verified: sc.d $t0, $t1, 8.
	off, err := New().Imm14(8)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x230009ac),
		ctorWord(t, New().ScD(lreg(t, 12), lreg(t, 13), off)),
	)

	in := New().ScD(lreg(t, 1), lreg(t, 2), off)
	_, ok := in.(ScD)
	require.True(t, ok, "type = %T, want ScD", in)
}

func TestScDDecodeEncode(t *testing.T) {
	in := decodeScD(0x230009ac, 0x90000000)

	scd, ok := in.(ScD)
	require.True(t, ok, "type = %T, want ScD", in)
	require.Equal(t, "sc.d $t0, $t1, 8", scd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), scd.Addr())
	require.Equal(t, 4, scd.Len())
	require.Equal(t, int64(8), scd.off.val)
	require.Equal(t, uint32(0x230009ac), ctorWord(t, scd))
}
