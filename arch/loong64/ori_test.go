package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestOriCtor(t *testing.T) {
	// llvm-mc-verified: ori $t0, $t1, 3855.
	v, err := NewUImm12(3855)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x03bc3dac),
		ctorWord(t, NewOri(lreg(t, 12), lreg(t, 13), v)),
	)
}

func TestOriDecodeEncode(t *testing.T) {
	x, ok := decodeOri(0x03bc3dac, 0x90000000).(Ori)
	require.True(t, ok, "type = %T, want Ori", x)

	require.Equal(t, "ori $t0, $t1, 3855", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(3855), x.imm.val)
	require.Equal(t, uint32(0x03bc3dac), ctorWord(t, x))
}
