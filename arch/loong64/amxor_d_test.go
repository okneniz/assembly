package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmxorDCtor(t *testing.T) {
	// llvm-mc-verified: amxor.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3864b5cc),
		ctorWord(t, New().AmxorD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmxorD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmxorD)
	require.True(t, ok, "type = %T, want AmxorD", in)
}

func TestAmxorDDecodeEncode(t *testing.T) {
	in := decodeAmxorD(0x3864b5cc, 0x90000000)

	amxord, ok := in.(AmxorD)
	require.True(t, ok, "type = %T, want AmxorD", in)
	require.Equal(t, "amxor.d $t0, $t1, $t2", amxord.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amxord.Addr())
	require.Equal(t, 4, amxord.Len())
	require.Equal(t, uint32(0x3864b5cc), ctorWord(t, amxord))
}
