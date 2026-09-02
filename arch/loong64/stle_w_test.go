package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStleWCtor(t *testing.T) {
	// llvm-mc-verified: stle.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387f39ac),
		ctorWord(t, New().StleW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().StleW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(StleW)
	require.True(t, ok, "type = %T, want StleW", in)
}

func TestStleWDecodeEncode(t *testing.T) {
	in := decodeStleW(0x387f39ac, 0x90000000)

	x, ok := in.(StleW)
	require.True(t, ok, "type = %T, want StleW", in)
	require.Equal(t, "stle.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x387f39ac), ctorWord(t, x))
}
