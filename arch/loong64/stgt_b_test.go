package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStgtBCtor(t *testing.T) {
	// llvm-mc-verified: stgt.b $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387c39ac),
		ctorWord(t, New().StgtB(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().StgtB(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(StgtB)
	require.True(t, ok, "type = %T, want StgtB", in)
}

func TestStgtBDecodeEncode(t *testing.T) {
	in := decodeStgtB(0x387c39ac, 0x90000000)

	x, ok := in.(StgtB)
	require.True(t, ok, "type = %T, want StgtB", in)
	require.Equal(t, "stgt.b $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x387c39ac), ctorWord(t, x))
}
