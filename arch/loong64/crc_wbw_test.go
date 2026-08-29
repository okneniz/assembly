package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCrcWBWCtor(t *testing.T) {
	// llvm-mc-verified: crc.w.b.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x002439ac),
		ctorWord(t, NewCrcWBW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewCrcWBW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(CrcWBW)
	require.True(t, ok, "type = %T, want CrcWBW", in)
}

func TestCrcWBWDecodeEncode(t *testing.T) {
	in := decodeCrcWBW(0x002439ac, 0x90000000)

	x, ok := in.(CrcWBW)
	require.True(t, ok, "type = %T, want CrcWBW", in)
	require.Equal(t, "crc.w.b.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x002439ac), ctorWord(t, x))
}
