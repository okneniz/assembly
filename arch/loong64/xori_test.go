package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestXoriCtor(t *testing.T) {
	// llvm-mc-verified: xori $t0, $t1, 3855.
	v, err := NewUImm12(3855)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x03fc3dac),
		ctorWord(t, NewXori(lreg(t, 12), lreg(t, 13), v)),
	)
}

func TestXoriDecodeEncode(t *testing.T) {
	x, ok := decodeXori(0x03fc3dac, 0x90000000).(Xori)
	require.True(t, ok, "type = %T, want Xori", x)

	require.Equal(t, "xori $t0, $t1, 3855", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(3855), x.imm.val)
	require.Equal(t, uint32(0x03fc3dac), ctorWord(t, x))
}
