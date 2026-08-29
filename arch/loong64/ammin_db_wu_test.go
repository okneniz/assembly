package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmminDbWuCtor(t *testing.T) {
	// llvm-mc-verified: ammin_db.wu $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387135cc),
		ctorWord(t, NewAmminDbWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmminDbWu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmminDbWu)
	require.True(t, ok, "type = %T, want AmminDbWu", in)
}

func TestAmminDbWuDecodeEncode(t *testing.T) {
	in := decodeAmminDbWu(0x387135cc, 0x90000000)

	ammindbwu, ok := in.(AmminDbWu)
	require.True(t, ok, "type = %T, want AmminDbWu", in)
	require.Equal(t, "ammin_db.wu $t0, $t1, $t2", ammindbwu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammindbwu.Addr())
	require.Equal(t, 4, ammindbwu.Len())
	require.Equal(t, uint32(0x387135cc), ctorWord(t, ammindbwu))
}
