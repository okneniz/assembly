package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmcasDbHCtor(t *testing.T) {
	// llvm-mc-verified: amcas_db.h $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385ab5cc),
		ctorWord(t, NewAmcasDbH(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmcasDbH(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmcasDbH)
	require.True(t, ok, "type = %T, want AmcasDbH", in)
}

func TestAmcasDbHDecodeEncode(t *testing.T) {
	in := decodeAmcasDbH(0x385ab5cc, 0x90000000)

	amcasdbh, ok := in.(AmcasDbH)
	require.True(t, ok, "type = %T, want AmcasDbH", in)
	require.Equal(t, "amcas_db.h $t0, $t1, $t2", amcasdbh.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amcasdbh.Addr())
	require.Equal(t, 4, amcasdbh.Len())
	require.Equal(t, uint32(0x385ab5cc), ctorWord(t, amcasdbh))
}
