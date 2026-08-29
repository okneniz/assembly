package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmmaxDbWuCtor(t *testing.T) {
	// llvm-mc-verified: ammax_db.wu $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387035cc),
		ctorWord(t, NewAmmaxDbWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmmaxDbWu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmmaxDbWu)
	require.True(t, ok, "type = %T, want AmmaxDbWu", in)
}

func TestAmmaxDbWuDecodeEncode(t *testing.T) {
	in := decodeAmmaxDbWu(0x387035cc, 0x90000000)

	ammaxdbwu, ok := in.(AmmaxDbWu)
	require.True(t, ok, "type = %T, want AmmaxDbWu", in)
	require.Equal(t, "ammax_db.wu $t0, $t1, $t2", ammaxdbwu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammaxdbwu.Addr())
	require.Equal(t, 4, ammaxdbwu.Len())
	require.Equal(t, uint32(0x387035cc), ctorWord(t, ammaxdbwu))
}
