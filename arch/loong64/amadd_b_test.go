package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmaddBCtor(t *testing.T) {
	// llvm-mc-verified: amadd.b $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385d35cc),
		ctorWord(t, New().AmaddB(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmaddB(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmaddB)
	require.True(t, ok, "type = %T, want AmaddB", in)
}

func TestAmaddBDecodeEncode(t *testing.T) {
	in := decodeAmaddB(0x385d35cc, 0x90000000)

	amaddb, ok := in.(AmaddB)
	require.True(t, ok, "type = %T, want AmaddB", in)
	require.Equal(t, "amadd.b $t0, $t1, $t2", amaddb.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amaddb.Addr())
	require.Equal(t, 4, amaddb.Len())
	require.Equal(t, uint32(0x385d35cc), ctorWord(t, amaddb))
}
