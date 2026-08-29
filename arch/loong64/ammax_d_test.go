package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmmaxDCtor(t *testing.T) {
	// llvm-mc-verified: ammax.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3865b5cc),
		ctorWord(t, NewAmmaxD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmmaxD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmmaxD)
	require.True(t, ok, "type = %T, want AmmaxD", in)
}

func TestAmmaxDDecodeEncode(t *testing.T) {
	in := decodeAmmaxD(0x3865b5cc, 0x90000000)

	ammaxd, ok := in.(AmmaxD)
	require.True(t, ok, "type = %T, want AmmaxD", in)
	require.Equal(t, "ammax.d $t0, $t1, $t2", ammaxd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammaxd.Addr())
	require.Equal(t, 4, ammaxd.Len())
	require.Equal(t, uint32(0x3865b5cc), ctorWord(t, ammaxd))
}
