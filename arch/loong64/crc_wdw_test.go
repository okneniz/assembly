package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCrcWDWCtor(t *testing.T) {
	// llvm-mc-verified: crc.w.d.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0025b9ac),
		ctorWord(t, New().CrcWDW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().CrcWDW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(CrcWDW)
	require.True(t, ok, "type = %T, want CrcWDW", in)
}

func TestCrcWDWDecodeEncode(t *testing.T) {
	in := decodeCrcWDW(0x0025b9ac, 0x90000000)

	x, ok := in.(CrcWDW)
	require.True(t, ok, "type = %T, want CrcWDW", in)
	require.Equal(t, "crc.w.d.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x0025b9ac), ctorWord(t, x))
}
