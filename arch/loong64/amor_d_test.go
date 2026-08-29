package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmorDCtor(t *testing.T) {
	// llvm-mc-verified: amor.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3863b5cc),
		ctorWord(t, NewAmorD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmorD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmorD)
	require.True(t, ok, "type = %T, want AmorD", in)
}

func TestAmorDDecodeEncode(t *testing.T) {
	in := decodeAmorD(0x3863b5cc, 0x90000000)

	amord, ok := in.(AmorD)
	require.True(t, ok, "type = %T, want AmorD", in)
	require.Equal(t, "amor.d $t0, $t1, $t2", amord.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amord.Addr())
	require.Equal(t, 4, amord.Len())
	require.Equal(t, uint32(0x3863b5cc), ctorWord(t, amord))
}
