package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmmaxDbDuCtor(t *testing.T) {
	// llvm-mc-verified: ammax_db.du $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3870b5cc),
		ctorWord(t, New().AmmaxDbDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmmaxDbDu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmmaxDbDu)
	require.True(t, ok, "type = %T, want AmmaxDbDu", in)
}

func TestAmmaxDbDuDecodeEncode(t *testing.T) {
	in := decodeAmmaxDbDu(0x3870b5cc, 0x90000000)

	ammaxdbdu, ok := in.(AmmaxDbDu)
	require.True(t, ok, "type = %T, want AmmaxDbDu", in)
	require.Equal(t, "ammax_db.du $t0, $t1, $t2", ammaxdbdu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammaxdbdu.Addr())
	require.Equal(t, 4, ammaxdbdu.Len())
	require.Equal(t, uint32(0x3870b5cc), ctorWord(t, ammaxdbdu))
}
