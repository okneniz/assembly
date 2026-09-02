package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmminDbDuCtor(t *testing.T) {
	// llvm-mc-verified: ammin_db.du $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3871b5cc),
		ctorWord(t, New().AmminDbDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmminDbDu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmminDbDu)
	require.True(t, ok, "type = %T, want AmminDbDu", in)
}

func TestAmminDbDuDecodeEncode(t *testing.T) {
	in := decodeAmminDbDu(0x3871b5cc, 0x90000000)

	ammindbdu, ok := in.(AmminDbDu)
	require.True(t, ok, "type = %T, want AmminDbDu", in)
	require.Equal(t, "ammin_db.du $t0, $t1, $t2", ammindbdu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammindbdu.Addr())
	require.Equal(t, 4, ammindbdu.Len())
	require.Equal(t, uint32(0x3871b5cc), ctorWord(t, ammindbdu))
}
