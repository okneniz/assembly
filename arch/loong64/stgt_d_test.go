package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStgtDCtor(t *testing.T) {
	// llvm-mc-verified: stgt.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387db9ac),
		ctorWord(t, NewStgtD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewStgtD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(StgtD)
	require.True(t, ok, "type = %T, want StgtD", in)
}

func TestStgtDDecodeEncode(t *testing.T) {
	in := decodeStgtD(0x387db9ac, 0x90000000)

	x, ok := in.(StgtD)
	require.True(t, ok, "type = %T, want StgtD", in)
	require.Equal(t, "stgt.d $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x387db9ac), ctorWord(t, x))
}
