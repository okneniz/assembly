package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmswapBCtor(t *testing.T) {
	// llvm-mc-verified: amswap.b $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385c35cc),
		ctorWord(t, New().AmswapB(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmswapB(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmswapB)
	require.True(t, ok, "type = %T, want AmswapB", in)
}

func TestAmswapBDecodeEncode(t *testing.T) {
	in := decodeAmswapB(0x385c35cc, 0x90000000)

	amswapb, ok := in.(AmswapB)
	require.True(t, ok, "type = %T, want AmswapB", in)
	require.Equal(t, "amswap.b $t0, $t1, $t2", amswapb.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amswapb.Addr())
	require.Equal(t, 4, amswapb.Len())
	require.Equal(t, uint32(0x385c35cc), ctorWord(t, amswapb))
}
