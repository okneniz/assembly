package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmxorDbWCtor(t *testing.T) {
	// llvm-mc-verified: amxor_db.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386d35cc),
		ctorWord(t, NewAmxorDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmxorDbW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmxorDbW)
	require.True(t, ok, "type = %T, want AmxorDbW", in)
}

func TestAmxorDbWDecodeEncode(t *testing.T) {
	in := decodeAmxorDbW(0x386d35cc, 0x90000000)

	amxordbw, ok := in.(AmxorDbW)
	require.True(t, ok, "type = %T, want AmxorDbW", in)
	require.Equal(t, "amxor_db.w $t0, $t1, $t2", amxordbw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amxordbw.Addr())
	require.Equal(t, 4, amxordbw.Len())
	require.Equal(t, uint32(0x386d35cc), ctorWord(t, amxordbw))
}
