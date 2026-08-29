package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmswapDbDCtor(t *testing.T) {
	// llvm-mc-verified: amswap_db.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3869b5cc),
		ctorWord(t, NewAmswapDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmswapDbD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmswapDbD)
	require.True(t, ok, "type = %T, want AmswapDbD", in)
}

func TestAmswapDbDDecodeEncode(t *testing.T) {
	in := decodeAmswapDbD(0x3869b5cc, 0x90000000)

	amswapdbd, ok := in.(AmswapDbD)
	require.True(t, ok, "type = %T, want AmswapDbD", in)
	require.Equal(t, "amswap_db.d $t0, $t1, $t2", amswapdbd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amswapdbd.Addr())
	require.Equal(t, 4, amswapdbd.Len())
	require.Equal(t, uint32(0x3869b5cc), ctorWord(t, amswapdbd))
}
