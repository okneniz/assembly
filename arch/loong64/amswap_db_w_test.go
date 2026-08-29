package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmswapDbWCtor(t *testing.T) {
	// llvm-mc-verified: amswap_db.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386935cc),
		ctorWord(t, NewAmswapDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmswapDbW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmswapDbW)
	require.True(t, ok, "type = %T, want AmswapDbW", in)
}

func TestAmswapDbWDecodeEncode(t *testing.T) {
	in := decodeAmswapDbW(0x386935cc, 0x90000000)

	amswapdbw, ok := in.(AmswapDbW)
	require.True(t, ok, "type = %T, want AmswapDbW", in)
	require.Equal(t, "amswap_db.w $t0, $t1, $t2", amswapdbw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amswapdbw.Addr())
	require.Equal(t, 4, amswapdbw.Len())
	require.Equal(t, uint32(0x386935cc), ctorWord(t, amswapdbw))
}
