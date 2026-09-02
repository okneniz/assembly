package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStleBCtor(t *testing.T) {
	// llvm-mc-verified: stle.b $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387e39ac),
		ctorWord(t, New().StleB(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().StleB(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(StleB)
	require.True(t, ok, "type = %T, want StleB", in)
}

func TestStleBDecodeEncode(t *testing.T) {
	in := decodeStleB(0x387e39ac, 0x90000000)

	x, ok := in.(StleB)
	require.True(t, ok, "type = %T, want StleB", in)
	require.Equal(t, "stle.b $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x387e39ac), ctorWord(t, x))
}
