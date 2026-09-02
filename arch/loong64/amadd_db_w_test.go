package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmaddDbWCtor(t *testing.T) {
	// llvm-mc-verified: amadd_db.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386a35cc),
		ctorWord(t, New().AmaddDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmaddDbW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmaddDbW)
	require.True(t, ok, "type = %T, want AmaddDbW", in)
}

func TestAmaddDbWDecodeEncode(t *testing.T) {
	in := decodeAmaddDbW(0x386a35cc, 0x90000000)

	amadddbw, ok := in.(AmaddDbW)
	require.True(t, ok, "type = %T, want AmaddDbW", in)
	require.Equal(t, "amadd_db.w $t0, $t1, $t2", amadddbw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amadddbw.Addr())
	require.Equal(t, 4, amadddbw.Len())
	require.Equal(t, uint32(0x386a35cc), ctorWord(t, amadddbw))
}
