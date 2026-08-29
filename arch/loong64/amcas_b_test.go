package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmcasBCtor(t *testing.T) {
	// llvm-mc-verified: amcas.b $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385835cc),
		ctorWord(t, NewAmcasB(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmcasB(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmcasB)
	require.True(t, ok, "type = %T, want AmcasB", in)
}

func TestAmcasBDecodeEncode(t *testing.T) {
	in := decodeAmcasB(0x385835cc, 0x90000000)

	amcasb, ok := in.(AmcasB)
	require.True(t, ok, "type = %T, want AmcasB", in)
	require.Equal(t, "amcas.b $t0, $t1, $t2", amcasb.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amcasb.Addr())
	require.Equal(t, 4, amcasb.Len())
	require.Equal(t, uint32(0x385835cc), ctorWord(t, amcasb))
}
