package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmcasDbWCtor(t *testing.T) {
	// llvm-mc-verified: amcas_db.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385b35cc),
		ctorWord(t, NewAmcasDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmcasDbW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmcasDbW)
	require.True(t, ok, "type = %T, want AmcasDbW", in)
}

func TestAmcasDbWDecodeEncode(t *testing.T) {
	in := decodeAmcasDbW(0x385b35cc, 0x90000000)

	amcasdbw, ok := in.(AmcasDbW)
	require.True(t, ok, "type = %T, want AmcasDbW", in)
	require.Equal(t, "amcas_db.w $t0, $t1, $t2", amcasdbw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amcasdbw.Addr())
	require.Equal(t, 4, amcasdbw.Len())
	require.Equal(t, uint32(0x385b35cc), ctorWord(t, amcasdbw))
}
