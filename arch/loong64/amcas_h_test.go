package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmcasHCtor(t *testing.T) {
	// llvm-mc-verified: amcas.h $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3858b5cc),
		ctorWord(t, New().AmcasH(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmcasH(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmcasH)
	require.True(t, ok, "type = %T, want AmcasH", in)
}

func TestAmcasHDecodeEncode(t *testing.T) {
	in := decodeAmcasH(0x3858b5cc, 0x90000000)

	amcash, ok := in.(AmcasH)
	require.True(t, ok, "type = %T, want AmcasH", in)
	require.Equal(t, "amcas.h $t0, $t1, $t2", amcash.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amcash.Addr())
	require.Equal(t, 4, amcash.Len())
	require.Equal(t, uint32(0x3858b5cc), ctorWord(t, amcash))
}
