package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmswapDbBCtor(t *testing.T) {
	// llvm-mc-verified: amswap_db.b $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385e35cc),
		ctorWord(t, NewAmswapDbB(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmswapDbB(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmswapDbB)
	require.True(t, ok, "type = %T, want AmswapDbB", in)
}

func TestAmswapDbBDecodeEncode(t *testing.T) {
	in := decodeAmswapDbB(0x385e35cc, 0x90000000)

	amswapdbb, ok := in.(AmswapDbB)
	require.True(t, ok, "type = %T, want AmswapDbB", in)
	require.Equal(t, "amswap_db.b $t0, $t1, $t2", amswapdbb.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amswapdbb.Addr())
	require.Equal(t, 4, amswapdbb.Len())
	require.Equal(t, uint32(0x385e35cc), ctorWord(t, amswapdbb))
}
