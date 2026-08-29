package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStgtWCtor(t *testing.T) {
	// llvm-mc-verified: stgt.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387d39ac),
		ctorWord(t, NewStgtW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewStgtW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(StgtW)
	require.True(t, ok, "type = %T, want StgtW", in)
}

func TestStgtWDecodeEncode(t *testing.T) {
	in := decodeStgtW(0x387d39ac, 0x90000000)

	x, ok := in.(StgtW)
	require.True(t, ok, "type = %T, want StgtW", in)
	require.Equal(t, "stgt.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x387d39ac), ctorWord(t, x))
}
