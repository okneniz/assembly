package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmandDbDCtor(t *testing.T) {
	// llvm-mc-verified: amand_db.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386bb5cc),
		ctorWord(t, NewAmandDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmandDbD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmandDbD)
	require.True(t, ok, "type = %T, want AmandDbD", in)
}

func TestAmandDbDDecodeEncode(t *testing.T) {
	in := decodeAmandDbD(0x386bb5cc, 0x90000000)

	amanddbd, ok := in.(AmandDbD)
	require.True(t, ok, "type = %T, want AmandDbD", in)
	require.Equal(t, "amand_db.d $t0, $t1, $t2", amanddbd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amanddbd.Addr())
	require.Equal(t, 4, amanddbd.Len())
	require.Equal(t, uint32(0x386bb5cc), ctorWord(t, amanddbd))
}
