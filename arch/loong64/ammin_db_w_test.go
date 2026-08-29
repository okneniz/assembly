package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmminDbWCtor(t *testing.T) {
	// llvm-mc-verified: ammin_db.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386f35cc),
		ctorWord(t, NewAmminDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmminDbW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmminDbW)
	require.True(t, ok, "type = %T, want AmminDbW", in)
}

func TestAmminDbWDecodeEncode(t *testing.T) {
	in := decodeAmminDbW(0x386f35cc, 0x90000000)

	ammindbw, ok := in.(AmminDbW)
	require.True(t, ok, "type = %T, want AmminDbW", in)
	require.Equal(t, "ammin_db.w $t0, $t1, $t2", ammindbw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammindbw.Addr())
	require.Equal(t, 4, ammindbw.Len())
	require.Equal(t, uint32(0x386f35cc), ctorWord(t, ammindbw))
}
