package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmswapDbHCtor(t *testing.T) {
	// llvm-mc-verified: amswap_db.h $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385eb5cc),
		ctorWord(t, NewAmswapDbH(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmswapDbH(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmswapDbH)
	require.True(t, ok, "type = %T, want AmswapDbH", in)
}

func TestAmswapDbHDecodeEncode(t *testing.T) {
	in := decodeAmswapDbH(0x385eb5cc, 0x90000000)

	amswapdbh, ok := in.(AmswapDbH)
	require.True(t, ok, "type = %T, want AmswapDbH", in)
	require.Equal(t, "amswap_db.h $t0, $t1, $t2", amswapdbh.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amswapdbh.Addr())
	require.Equal(t, 4, amswapdbh.Len())
	require.Equal(t, uint32(0x385eb5cc), ctorWord(t, amswapdbh))
}
