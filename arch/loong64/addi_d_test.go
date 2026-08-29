package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAddiDCtor(t *testing.T) {
	// llvm-mc-verified: addi.d $t0, $t1, -16.
	require.Equal(
		t,
		uint32(0x02ffc1ac),
		ctorWord(t, NewAddiD(lreg(t, 12), lreg(t, 13), imm12v(t, -16))),
	)
}

func TestAddiDDecodeEncode(t *testing.T) {
	x, ok := decodeAddiD(0x02ffc1ac, 0x90000000).(AddiD)
	require.True(t, ok, "type = %T, want AddiD", x)

	require.Equal(t, "addi.d $t0, $t1, -16", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(-16), x.imm.val)

	// The negative immediate round-trips through the sign-extended field.
	require.Equal(t, uint32(0x02ffc1ac), ctorWord(t, x))
}
