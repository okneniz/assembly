package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmcasDbBCtor(t *testing.T) {
	// llvm-mc-verified: amcas_db.b $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385a35cc),
		ctorWord(t, NewAmcasDbB(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmcasDbB(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmcasDbB)
	require.True(t, ok, "type = %T, want AmcasDbB", in)
}

func TestAmcasDbBDecodeEncode(t *testing.T) {
	in := decodeAmcasDbB(0x385a35cc, 0x90000000)

	amcasdbb, ok := in.(AmcasDbB)
	require.True(t, ok, "type = %T, want AmcasDbB", in)
	require.Equal(t, "amcas_db.b $t0, $t1, $t2", amcasdbb.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amcasdbb.Addr())
	require.Equal(t, 4, amcasdbb.Len())
	require.Equal(t, uint32(0x385a35cc), ctorWord(t, amcasdbb))
}
