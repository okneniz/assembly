package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmmaxDbDCtor(t *testing.T) {
	// llvm-mc-verified: ammax_db.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386eb5cc),
		ctorWord(t, NewAmmaxDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmmaxDbD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmmaxDbD)
	require.True(t, ok, "type = %T, want AmmaxDbD", in)
}

func TestAmmaxDbDDecodeEncode(t *testing.T) {
	in := decodeAmmaxDbD(0x386eb5cc, 0x90000000)

	ammaxdbd, ok := in.(AmmaxDbD)
	require.True(t, ok, "type = %T, want AmmaxDbD", in)
	require.Equal(t, "ammax_db.d $t0, $t1, $t2", ammaxdbd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammaxdbd.Addr())
	require.Equal(t, 4, ammaxdbd.Len())
	require.Equal(t, uint32(0x386eb5cc), ctorWord(t, ammaxdbd))
}
