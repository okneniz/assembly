package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmaddDbHCtor(t *testing.T) {
	// llvm-mc-verified: amadd_db.h $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385fb5cc),
		ctorWord(t, New().AmaddDbH(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmaddDbH(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmaddDbH)
	require.True(t, ok, "type = %T, want AmaddDbH", in)
}

func TestAmaddDbHDecodeEncode(t *testing.T) {
	in := decodeAmaddDbH(0x385fb5cc, 0x90000000)

	amadddbh, ok := in.(AmaddDbH)
	require.True(t, ok, "type = %T, want AmaddDbH", in)
	require.Equal(t, "amadd_db.h $t0, $t1, $t2", amadddbh.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amadddbh.Addr())
	require.Equal(t, 4, amadddbh.Len())
	require.Equal(t, uint32(0x385fb5cc), ctorWord(t, amadddbh))
}
