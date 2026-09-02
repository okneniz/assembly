package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmswapHCtor(t *testing.T) {
	// llvm-mc-verified: amswap.h $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385cb5cc),
		ctorWord(t, New().AmswapH(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmswapH(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmswapH)
	require.True(t, ok, "type = %T, want AmswapH", in)
}

func TestAmswapHDecodeEncode(t *testing.T) {
	in := decodeAmswapH(0x385cb5cc, 0x90000000)

	amswaph, ok := in.(AmswapH)
	require.True(t, ok, "type = %T, want AmswapH", in)
	require.Equal(t, "amswap.h $t0, $t1, $t2", amswaph.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amswaph.Addr())
	require.Equal(t, 4, amswaph.Len())
	require.Equal(t, uint32(0x385cb5cc), ctorWord(t, amswaph))
}
