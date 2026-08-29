package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRotriWCtor(t *testing.T) {
	// llvm-mc-verified: rotri.w $t0, $t1, 3 (the shift amount is the
	// ui5 field [14:10]; bit 15 belongs to the opcode).
	v, err := NewUImm5(3)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x004c8dac),
		ctorWord(t, NewRotriW(lreg(t, 12), lreg(t, 13), v)),
	)
}

func TestRotriWDecodeEncode(t *testing.T) {
	x, ok := decodeRotriW(0x004c8dac, 0x90000000).(RotriW)
	require.True(t, ok, "type = %T, want RotriW", x)

	require.Equal(t, "rotri.w $t0, $t1, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(3), x.imm.val)
	require.Equal(t, uint32(0x004c8dac), ctorWord(t, x))
}
