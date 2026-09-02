package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSraiDCtor(t *testing.T) {
	// llvm-mc-verified: srai.d $t0, $t1, 3.
	v, err := New().UImm6(3)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x00490dac),
		ctorWord(t, New().SraiD(lreg(t, 12), lreg(t, 13), v)),
	)
}

func TestSraiDDecodeEncode(t *testing.T) {
	x, ok := decodeSraiD(0x00490dac, 0x90000000).(SraiD)
	require.True(t, ok, "type = %T, want SraiD", x)

	require.Equal(t, "srai.d $t0, $t1, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(3), x.imm.val)
	require.Equal(t, uint32(0x00490dac), ctorWord(t, x))
}
