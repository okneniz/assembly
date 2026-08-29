package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSlliDCtor(t *testing.T) {
	// llvm-mc-verified: slli.d $t0, $t1, 3.
	v, err := NewUImm6(3)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x00410dac),
		ctorWord(t, NewSlliD(lreg(t, 12), lreg(t, 13), v)),
	)
}

func TestSlliDDecodeEncode(t *testing.T) {
	x, ok := decodeSlliD(0x00410dac, 0x90000000).(SlliD)
	require.True(t, ok, "type = %T, want SlliD", x)

	require.Equal(t, "slli.d $t0, $t1, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(3), x.imm.val)
	require.Equal(t, uint32(0x00410dac), ctorWord(t, x))
}
