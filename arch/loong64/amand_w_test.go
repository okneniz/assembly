package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmandWCtor(t *testing.T) {
	// llvm-mc-verified: amand.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386235cc),
		ctorWord(t, New().AmandW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmandW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmandW)
	require.True(t, ok, "type = %T, want AmandW", in)
}

func TestAmandWDecodeEncode(t *testing.T) {
	in := decodeAmandW(0x386235cc, 0x90000000)

	amandw, ok := in.(AmandW)
	require.True(t, ok, "type = %T, want AmandW", in)
	require.Equal(t, "amand.w $t0, $t1, $t2", amandw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amandw.Addr())
	require.Equal(t, 4, amandw.Len())
	require.Equal(t, uint32(0x386235cc), ctorWord(t, amandw))
}
