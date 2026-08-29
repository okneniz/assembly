package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCtzDCtor(t *testing.T) {
	// llvm-mc-verified: ctz.d $t0, $t1.
	require.Equal(
		t,
		uint32(0x00002dac),
		ctorWord(t, NewCtzD(lreg(t, 12), lreg(t, 13))),
	)

	in := NewCtzD(lreg(t, 1), lreg(t, 2))
	_, ok := in.(CtzD)
	require.True(t, ok, "type = %T, want CtzD", in)
}

func TestCtzDDecodeEncode(t *testing.T) {
	in := decodeCtzD(0x00002dac, 0x90000000)

	ctzd, ok := in.(CtzD)
	require.True(t, ok, "type = %T, want CtzD", in)
	require.Equal(t, "ctz.d $t0, $t1", ctzd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ctzd.Addr())
	require.Equal(t, uint32(0x00002dac), ctorWord(t, ctzd))
}
