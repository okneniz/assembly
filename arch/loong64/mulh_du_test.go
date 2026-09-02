package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestMulhDuCtor(t *testing.T) {
	// llvm-mc-verified: mulh.du $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001eb9ac),
		ctorWord(t, New().MulhDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().MulhDu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(MulhDu)
	require.True(t, ok, "type = %T, want MulhDu", in)
}

func TestMulhDuDecodeEncode(t *testing.T) {
	in := decodeMulhDu(0x001eb9ac, 0x90000000)

	mulhdu, ok := in.(MulhDu)
	require.True(t, ok, "type = %T, want MulhDu", in)
	require.Equal(t, "mulh.du $t0, $t1, $t2", mulhdu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), mulhdu.Addr())
	require.Equal(t, uint32(0x001eb9ac), ctorWord(t, mulhdu))
}
