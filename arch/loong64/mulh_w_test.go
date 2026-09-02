package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestMulhWCtor(t *testing.T) {
	// llvm-mc-verified: mulh.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001cb9ac),
		ctorWord(t, New().MulhW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().MulhW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(MulhW)
	require.True(t, ok, "type = %T, want MulhW", in)
}

func TestMulhWDecodeEncode(t *testing.T) {
	in := decodeMulhW(0x001cb9ac, 0x90000000)

	mulhw, ok := in.(MulhW)
	require.True(t, ok, "type = %T, want MulhW", in)
	require.Equal(t, "mulh.w $t0, $t1, $t2", mulhw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), mulhw.Addr())
	require.Equal(t, uint32(0x001cb9ac), ctorWord(t, mulhw))
}
