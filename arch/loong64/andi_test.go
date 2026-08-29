package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAndiCtor(t *testing.T) {
	// llvm-mc-verified: andi $t0, $t1, 3855.
	v, err := NewUImm12(3855)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x037c3dac),
		ctorWord(t, NewAndi(lreg(t, 12), lreg(t, 13), v)),
	)
}

func TestAndiDecodeEncode(t *testing.T) {
	x, ok := decodeAndi(0x037c3dac, 0x90000000).(Andi)
	require.True(t, ok, "type = %T, want Andi", x)

	require.Equal(t, "andi $t0, $t1, 3855", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(3855), x.imm.val)
	require.Equal(t, uint32(0x037c3dac), ctorWord(t, x))
}
