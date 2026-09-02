package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStxBCtor(t *testing.T) {
	// llvm-mc-verified: stx.b $t0, $t1, $t2.
	in := New().StxB(lreg(t, 12), lreg(t, 13), lreg(t, 14))
	require.Equal(t, uint32(0x381039ac), ctorWord(t, in))

	_, ok := in.(StxB)
	require.True(t, ok, "type = %T, want StxB", in)
}

func TestStxBDecodeEncode(t *testing.T) {
	in := decodeStxB(0x381039ac, 0x90000000)

	x, ok := in.(StxB)
	require.True(t, ok, "type = %T, want StxB", in)
	require.Equal(t, "stx.b $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x381039ac), ctorWord(t, x))
}
