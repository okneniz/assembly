package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSraiWCtor(t *testing.T) {
	// llvm-mc-verified: srai.w $t0, $t1, 3.
	v, err := New().UImm5(3)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x00488dac),
		ctorWord(t, New().SraiW(lreg(t, 12), lreg(t, 13), v)),
	)
}

func TestSraiWDecodeEncode(t *testing.T) {
	x, ok := decodeSraiW(0x00488dac, 0x90000000).(SraiW)
	require.True(t, ok, "type = %T, want SraiW", x)

	require.Equal(t, "srai.w $t0, $t1, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(3), x.imm.val)
	require.Equal(t, uint32(0x00488dac), ctorWord(t, x))
}
