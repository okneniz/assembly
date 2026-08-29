package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmcasDbDCtor(t *testing.T) {
	// llvm-mc-verified: amcas_db.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385bb5cc),
		ctorWord(t, NewAmcasDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmcasDbD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmcasDbD)
	require.True(t, ok, "type = %T, want AmcasDbD", in)
}

func TestAmcasDbDDecodeEncode(t *testing.T) {
	in := decodeAmcasDbD(0x385bb5cc, 0x90000000)

	amcasdbd, ok := in.(AmcasDbD)
	require.True(t, ok, "type = %T, want AmcasDbD", in)
	require.Equal(t, "amcas_db.d $t0, $t1, $t2", amcasdbd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amcasdbd.Addr())
	require.Equal(t, 4, amcasdbd.Len())
	require.Equal(t, uint32(0x385bb5cc), ctorWord(t, amcasdbd))
}
