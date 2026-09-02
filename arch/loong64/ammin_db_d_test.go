package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmminDbDCtor(t *testing.T) {
	// llvm-mc-verified: ammin_db.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386fb5cc),
		ctorWord(t, New().AmminDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmminDbD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmminDbD)
	require.True(t, ok, "type = %T, want AmminDbD", in)
}

func TestAmminDbDDecodeEncode(t *testing.T) {
	in := decodeAmminDbD(0x386fb5cc, 0x90000000)

	ammindbd, ok := in.(AmminDbD)
	require.True(t, ok, "type = %T, want AmminDbD", in)
	require.Equal(t, "ammin_db.d $t0, $t1, $t2", ammindbd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammindbd.Addr())
	require.Equal(t, 4, ammindbd.Len())
	require.Equal(t, uint32(0x386fb5cc), ctorWord(t, ammindbd))
}
