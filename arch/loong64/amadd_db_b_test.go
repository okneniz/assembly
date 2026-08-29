package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmaddDbBCtor(t *testing.T) {
	// llvm-mc-verified: amadd_db.b $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385f35cc),
		ctorWord(t, NewAmaddDbB(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmaddDbB(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmaddDbB)
	require.True(t, ok, "type = %T, want AmaddDbB", in)
}

func TestAmaddDbBDecodeEncode(t *testing.T) {
	in := decodeAmaddDbB(0x385f35cc, 0x90000000)

	amadddbb, ok := in.(AmaddDbB)
	require.True(t, ok, "type = %T, want AmaddDbB", in)
	require.Equal(t, "amadd_db.b $t0, $t1, $t2", amadddbb.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amadddbb.Addr())
	require.Equal(t, 4, amadddbb.Len())
	require.Equal(t, uint32(0x385f35cc), ctorWord(t, amadddbb))
}
