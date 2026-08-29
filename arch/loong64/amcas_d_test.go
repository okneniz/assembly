package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmcasDCtor(t *testing.T) {
	// llvm-mc-verified: amcas.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3859b5cc),
		ctorWord(t, NewAmcasD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmcasD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmcasD)
	require.True(t, ok, "type = %T, want AmcasD", in)
}

func TestAmcasDDecodeEncode(t *testing.T) {
	in := decodeAmcasD(0x3859b5cc, 0x90000000)

	amcasd, ok := in.(AmcasD)
	require.True(t, ok, "type = %T, want AmcasD", in)
	require.Equal(t, "amcas.d $t0, $t1, $t2", amcasd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amcasd.Addr())
	require.Equal(t, 4, amcasd.Len())
	require.Equal(t, uint32(0x3859b5cc), ctorWord(t, amcasd))
}
