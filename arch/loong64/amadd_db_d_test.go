package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmaddDbDCtor(t *testing.T) {
	// llvm-mc-verified: amadd_db.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386ab5cc),
		ctorWord(t, New().AmaddDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmaddDbD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmaddDbD)
	require.True(t, ok, "type = %T, want AmaddDbD", in)
}

func TestAmaddDbDDecodeEncode(t *testing.T) {
	in := decodeAmaddDbD(0x386ab5cc, 0x90000000)

	amadddbd, ok := in.(AmaddDbD)
	require.True(t, ok, "type = %T, want AmaddDbD", in)
	require.Equal(t, "amadd_db.d $t0, $t1, $t2", amadddbd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amadddbd.Addr())
	require.Equal(t, 4, amadddbd.Len())
	require.Equal(t, uint32(0x386ab5cc), ctorWord(t, amadddbd))
}
