package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmxorWCtor(t *testing.T) {
	// llvm-mc-verified: amxor.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386435cc),
		ctorWord(t, New().AmxorW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmxorW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmxorW)
	require.True(t, ok, "type = %T, want AmxorW", in)
}

func TestAmxorWDecodeEncode(t *testing.T) {
	in := decodeAmxorW(0x386435cc, 0x90000000)

	amxorw, ok := in.(AmxorW)
	require.True(t, ok, "type = %T, want AmxorW", in)
	require.Equal(t, "amxor.w $t0, $t1, $t2", amxorw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amxorw.Addr())
	require.Equal(t, 4, amxorw.Len())
	require.Equal(t, uint32(0x386435cc), ctorWord(t, amxorw))
}
