package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmandDCtor(t *testing.T) {
	// llvm-mc-verified: amand.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3862b5cc),
		ctorWord(t, New().AmandD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmandD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmandD)
	require.True(t, ok, "type = %T, want AmandD", in)
}

func TestAmandDDecodeEncode(t *testing.T) {
	in := decodeAmandD(0x3862b5cc, 0x90000000)

	amandd, ok := in.(AmandD)
	require.True(t, ok, "type = %T, want AmandD", in)
	require.Equal(t, "amand.d $t0, $t1, $t2", amandd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amandd.Addr())
	require.Equal(t, 4, amandd.Len())
	require.Equal(t, uint32(0x3862b5cc), ctorWord(t, amandd))
}
