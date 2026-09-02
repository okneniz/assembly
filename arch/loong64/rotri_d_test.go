package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRotriDCtor(t *testing.T) {
	// llvm-mc-verified: rotri.d $t0, $t1, 3.
	v, err := New().UImm6(3)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x004d0dac),
		ctorWord(t, New().RotriD(lreg(t, 12), lreg(t, 13), v)),
	)
}

func TestRotriDDecodeEncode(t *testing.T) {
	x, ok := decodeRotriD(0x004d0dac, 0x90000000).(RotriD)
	require.True(t, ok, "type = %T, want RotriD", x)

	require.Equal(t, "rotri.d $t0, $t1, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(3), x.imm.val)
	require.Equal(t, uint32(0x004d0dac), ctorWord(t, x))
}
