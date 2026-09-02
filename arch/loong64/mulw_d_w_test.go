package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestMulwDWCtor(t *testing.T) {
	// llvm-mc-verified: mulw.d.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001f39ac),
		ctorWord(t, New().MulwDW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().MulwDW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(MulwDW)
	require.True(t, ok, "type = %T, want MulwDW", in)
}

func TestMulwDWDecodeEncode(t *testing.T) {
	in := decodeMulwDW(0x001f39ac, 0x90000000)

	mulwdw, ok := in.(MulwDW)
	require.True(t, ok, "type = %T, want MulwDW", in)
	require.Equal(t, "mulw.d.w $t0, $t1, $t2", mulwdw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), mulwdw.Addr())
	require.Equal(t, uint32(0x001f39ac), ctorWord(t, mulwdw))
}
