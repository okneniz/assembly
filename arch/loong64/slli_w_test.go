package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSlliWCtor(t *testing.T) {
	// llvm-mc-verified: slli.w $t0, $t1, 3.
	v, err := NewUImm5(3)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x00408dac),
		ctorWord(t, NewSlliW(lreg(t, 12), lreg(t, 13), v)),
	)
}

func TestSlliWDecodeEncode(t *testing.T) {
	x, ok := decodeSlliW(0x00408dac, 0x90000000).(SlliW)
	require.True(t, ok, "type = %T, want SlliW", x)

	require.Equal(t, "slli.w $t0, $t1, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(3), x.imm.val)
	require.Equal(t, uint32(0x00408dac), ctorWord(t, x))
}
