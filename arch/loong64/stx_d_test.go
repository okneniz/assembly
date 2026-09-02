package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStxDCtor(t *testing.T) {
	// llvm-mc-verified: stx.d $t0, $t1, $t2.
	in := New().StxD(lreg(t, 12), lreg(t, 13), lreg(t, 14))
	require.Equal(t, uint32(0x381c39ac), ctorWord(t, in))

	_, ok := in.(StxD)
	require.True(t, ok, "type = %T, want StxD", in)
}

func TestStxDDecodeEncode(t *testing.T) {
	in := decodeStxD(0x381c39ac, 0x90000000)

	x, ok := in.(StxD)
	require.True(t, ok, "type = %T, want StxD", in)
	require.Equal(t, "stx.d $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x381c39ac), ctorWord(t, x))
}
