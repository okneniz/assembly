package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCrcWWWCtor(t *testing.T) {
	// llvm-mc-verified: crc.w.w.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x002539ac),
		ctorWord(t, NewCrcWWW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewCrcWWW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(CrcWWW)
	require.True(t, ok, "type = %T, want CrcWWW", in)
}

func TestCrcWWWDecodeEncode(t *testing.T) {
	in := decodeCrcWWW(0x002539ac, 0x90000000)

	x, ok := in.(CrcWWW)
	require.True(t, ok, "type = %T, want CrcWWW", in)
	require.Equal(t, "crc.w.w.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x002539ac), ctorWord(t, x))
}
