package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStxHCtor(t *testing.T) {
	// llvm-mc-verified: stx.h $t0, $t1, $t2.
	in := New().StxH(lreg(t, 12), lreg(t, 13), lreg(t, 14))
	require.Equal(t, uint32(0x381439ac), ctorWord(t, in))

	_, ok := in.(StxH)
	require.True(t, ok, "type = %T, want StxH", in)
}

func TestStxHDecodeEncode(t *testing.T) {
	in := decodeStxH(0x381439ac, 0x90000000)

	x, ok := in.(StxH)
	require.True(t, ok, "type = %T, want StxH", in)
	require.Equal(t, "stx.h $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x381439ac), ctorWord(t, x))
}
