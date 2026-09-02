package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestMulWCtor(t *testing.T) {
	// llvm-mc-verified: mul.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001c39ac),
		ctorWord(t, New().MulW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().MulW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(MulW)
	require.True(t, ok, "type = %T, want MulW", in)
}

func TestMulWDecodeEncode(t *testing.T) {
	in := decodeMulW(0x001c39ac, 0x90000000)

	mulw, ok := in.(MulW)
	require.True(t, ok, "type = %T, want MulW", in)
	require.Equal(t, "mul.w $t0, $t1, $t2", mulw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), mulw.Addr())
	require.Equal(t, uint32(0x001c39ac), ctorWord(t, mulw))
}
