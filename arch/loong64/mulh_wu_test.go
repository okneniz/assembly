package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestMulhWuCtor(t *testing.T) {
	// llvm-mc-verified: mulh.wu $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001d39ac),
		ctorWord(t, New().MulhWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().MulhWu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(MulhWu)
	require.True(t, ok, "type = %T, want MulhWu", in)
}

func TestMulhWuDecodeEncode(t *testing.T) {
	in := decodeMulhWu(0x001d39ac, 0x90000000)

	mulhwu, ok := in.(MulhWu)
	require.True(t, ok, "type = %T, want MulhWu", in)
	require.Equal(t, "mulh.wu $t0, $t1, $t2", mulhwu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), mulhwu.Addr())
	require.Equal(t, uint32(0x001d39ac), ctorWord(t, mulhwu))
}
