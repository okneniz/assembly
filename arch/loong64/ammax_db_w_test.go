package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmmaxDbWCtor(t *testing.T) {
	// llvm-mc-verified: ammax_db.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386e35cc),
		ctorWord(t, New().AmmaxDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmmaxDbW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmmaxDbW)
	require.True(t, ok, "type = %T, want AmmaxDbW", in)
}

func TestAmmaxDbWDecodeEncode(t *testing.T) {
	in := decodeAmmaxDbW(0x386e35cc, 0x90000000)

	ammaxdbw, ok := in.(AmmaxDbW)
	require.True(t, ok, "type = %T, want AmmaxDbW", in)
	require.Equal(t, "ammax_db.w $t0, $t1, $t2", ammaxdbw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammaxdbw.Addr())
	require.Equal(t, 4, ammaxdbw.Len())
	require.Equal(t, uint32(0x386e35cc), ctorWord(t, ammaxdbw))
}
