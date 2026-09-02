package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestMulDCtor(t *testing.T) {
	// llvm-mc-verified: mul.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001db9ac),
		ctorWord(t, New().MulD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().MulD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(MulD)
	require.True(t, ok, "type = %T, want MulD", in)
}

func TestMulDDecodeEncode(t *testing.T) {
	in := decodeMulD(0x001db9ac, 0x90000000)

	muld, ok := in.(MulD)
	require.True(t, ok, "type = %T, want MulD", in)
	require.Equal(t, "mul.d $t0, $t1, $t2", muld.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), muld.Addr())
	require.Equal(t, uint32(0x001db9ac), ctorWord(t, muld))
}
