package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSltiCtor(t *testing.T) {
	// llvm-mc-verified: slti $t0, $t1, -16.
	require.Equal(
		t,
		uint32(0x023fc1ac),
		ctorWord(t, NewSlti(lreg(t, 12), lreg(t, 13), imm12v(t, -16))),
	)
}

func TestSltiDecodeEncode(t *testing.T) {
	x, ok := decodeSlti(0x023fc1ac, 0x90000000).(Slti)
	require.True(t, ok, "type = %T, want Slti", x)

	require.Equal(t, "slti $t0, $t1, -16", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(-16), x.imm.val)

	// The negative immediate round-trips through the sign-extended field.
	require.Equal(t, uint32(0x023fc1ac), ctorWord(t, x))
}
