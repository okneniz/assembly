package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestClzDCtor(t *testing.T) {
	// llvm-mc-verified: clz.d $t0, $t1.
	require.Equal(
		t,
		uint32(0x000025ac),
		ctorWord(t, New().ClzD(lreg(t, 12), lreg(t, 13))),
	)

	in := New().ClzD(lreg(t, 1), lreg(t, 2))
	_, ok := in.(ClzD)
	require.True(t, ok, "type = %T, want ClzD", in)
}

func TestClzDDecodeEncode(t *testing.T) {
	in := decodeClzD(0x000025ac, 0x90000000)

	clzd, ok := in.(ClzD)
	require.True(t, ok, "type = %T, want ClzD", in)
	require.Equal(t, "clz.d $t0, $t1", clzd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), clzd.Addr())
	require.Equal(t, uint32(0x000025ac), ctorWord(t, clzd))
}
