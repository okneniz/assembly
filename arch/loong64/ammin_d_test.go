package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmminDCtor(t *testing.T) {
	// llvm-mc-verified: ammin.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3866b5cc),
		ctorWord(t, NewAmminD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmminD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmminD)
	require.True(t, ok, "type = %T, want AmminD", in)
}

func TestAmminDDecodeEncode(t *testing.T) {
	in := decodeAmminD(0x3866b5cc, 0x90000000)

	ammind, ok := in.(AmminD)
	require.True(t, ok, "type = %T, want AmminD", in)
	require.Equal(t, "ammin.d $t0, $t1, $t2", ammind.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammind.Addr())
	require.Equal(t, 4, ammind.Len())
	require.Equal(t, uint32(0x3866b5cc), ctorWord(t, ammind))
}
