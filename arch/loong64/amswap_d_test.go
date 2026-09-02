package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmswapDCtor(t *testing.T) {
	// llvm-mc-verified: amswap.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3860b5cc),
		ctorWord(t, New().AmswapD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmswapD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmswapD)
	require.True(t, ok, "type = %T, want AmswapD", in)
}

func TestAmswapDDecodeEncode(t *testing.T) {
	in := decodeAmswapD(0x3860b5cc, 0x90000000)

	amswapd, ok := in.(AmswapD)
	require.True(t, ok, "type = %T, want AmswapD", in)
	require.Equal(t, "amswap.d $t0, $t1, $t2", amswapd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amswapd.Addr())
	require.Equal(t, 4, amswapd.Len())
	require.Equal(t, uint32(0x3860b5cc), ctorWord(t, amswapd))
}
