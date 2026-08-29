package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmandDbWCtor(t *testing.T) {
	// llvm-mc-verified: amand_db.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386b35cc),
		ctorWord(t, NewAmandDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmandDbW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmandDbW)
	require.True(t, ok, "type = %T, want AmandDbW", in)
}

func TestAmandDbWDecodeEncode(t *testing.T) {
	in := decodeAmandDbW(0x386b35cc, 0x90000000)

	amanddbw, ok := in.(AmandDbW)
	require.True(t, ok, "type = %T, want AmandDbW", in)
	require.Equal(t, "amand_db.w $t0, $t1, $t2", amanddbw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amanddbw.Addr())
	require.Equal(t, 4, amanddbw.Len())
	require.Equal(t, uint32(0x386b35cc), ctorWord(t, amanddbw))
}
