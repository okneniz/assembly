package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestMulhDCtor(t *testing.T) {
	// llvm-mc-verified: mulh.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001e39ac),
		ctorWord(t, NewMulhD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewMulhD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(MulhD)
	require.True(t, ok, "type = %T, want MulhD", in)
}

func TestMulhDDecodeEncode(t *testing.T) {
	in := decodeMulhD(0x001e39ac, 0x90000000)

	mulhd, ok := in.(MulhD)
	require.True(t, ok, "type = %T, want MulhD", in)
	require.Equal(t, "mulh.d $t0, $t1, $t2", mulhd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), mulhd.Addr())
	require.Equal(t, uint32(0x001e39ac), ctorWord(t, mulhd))
}
