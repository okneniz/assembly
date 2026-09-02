package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmmaxDuCtor(t *testing.T) {
	// llvm-mc-verified: ammax.du $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3867b5cc),
		ctorWord(t, New().AmmaxDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmmaxDu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmmaxDu)
	require.True(t, ok, "type = %T, want AmmaxDu", in)
}

func TestAmmaxDuDecodeEncode(t *testing.T) {
	in := decodeAmmaxDu(0x3867b5cc, 0x90000000)

	ammaxdu, ok := in.(AmmaxDu)
	require.True(t, ok, "type = %T, want AmmaxDu", in)
	require.Equal(t, "ammax.du $t0, $t1, $t2", ammaxdu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammaxdu.Addr())
	require.Equal(t, 4, ammaxdu.Len())
	require.Equal(t, uint32(0x3867b5cc), ctorWord(t, ammaxdu))
}
