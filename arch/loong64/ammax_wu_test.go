package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmmaxWuCtor(t *testing.T) {
	// llvm-mc-verified: ammax.wu $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386735cc),
		ctorWord(t, New().AmmaxWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmmaxWu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmmaxWu)
	require.True(t, ok, "type = %T, want AmmaxWu", in)
}

func TestAmmaxWuDecodeEncode(t *testing.T) {
	in := decodeAmmaxWu(0x386735cc, 0x90000000)

	ammaxwu, ok := in.(AmmaxWu)
	require.True(t, ok, "type = %T, want AmmaxWu", in)
	require.Equal(t, "ammax.wu $t0, $t1, $t2", ammaxwu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammaxwu.Addr())
	require.Equal(t, 4, ammaxwu.Len())
	require.Equal(t, uint32(0x386735cc), ctorWord(t, ammaxwu))
}
