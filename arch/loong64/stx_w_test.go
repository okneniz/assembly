package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStxWCtor(t *testing.T) {
	// llvm-mc-verified: stx.w $t0, $t1, $t2.
	in := New().StxW(lreg(t, 12), lreg(t, 13), lreg(t, 14))
	require.Equal(t, uint32(0x381839ac), ctorWord(t, in))

	_, ok := in.(StxW)
	require.True(t, ok, "type = %T, want StxW", in)
}

func TestStxWDecodeEncode(t *testing.T) {
	in := decodeStxW(0x381839ac, 0x90000000)

	x, ok := in.(StxW)
	require.True(t, ok, "type = %T, want StxW", in)
	require.Equal(t, "stx.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x381839ac), ctorWord(t, x))
}
