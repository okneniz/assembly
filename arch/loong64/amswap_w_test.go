package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmswapWCtor(t *testing.T) {
	// llvm-mc-verified: amswap.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386035cc),
		ctorWord(t, NewAmswapW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmswapW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmswapW)
	require.True(t, ok, "type = %T, want AmswapW", in)
}

func TestAmswapWDecodeEncode(t *testing.T) {
	in := decodeAmswapW(0x386035cc, 0x90000000)

	amswapw, ok := in.(AmswapW)
	require.True(t, ok, "type = %T, want AmswapW", in)
	require.Equal(t, "amswap.w $t0, $t1, $t2", amswapw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amswapw.Addr())
	require.Equal(t, 4, amswapw.Len())
	require.Equal(t, uint32(0x386035cc), ctorWord(t, amswapw))
}
