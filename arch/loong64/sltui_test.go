package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSltuiCtor(t *testing.T) {
	// llvm-mc-verified: sltui $t0, $t1, -16.
	require.Equal(
		t,
		uint32(0x027fc1ac),
		ctorWord(t, NewSltui(lreg(t, 12), lreg(t, 13), imm12v(t, -16))),
	)
}

func TestSltuiDecodeEncode(t *testing.T) {
	x, ok := decodeSltui(0x027fc1ac, 0x90000000).(Sltui)
	require.True(t, ok, "type = %T, want Sltui", x)

	require.Equal(t, "sltui $t0, $t1, -16", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(-16), x.imm.val)

	// The negative immediate round-trips through the sign-extended field.
	require.Equal(t, uint32(0x027fc1ac), ctorWord(t, x))
}
