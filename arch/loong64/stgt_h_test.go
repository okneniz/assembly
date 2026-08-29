package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStgtHCtor(t *testing.T) {
	// llvm-mc-verified: stgt.h $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387cb9ac),
		ctorWord(t, NewStgtH(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewStgtH(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(StgtH)
	require.True(t, ok, "type = %T, want StgtH", in)
}

func TestStgtHDecodeEncode(t *testing.T) {
	in := decodeStgtH(0x387cb9ac, 0x90000000)

	x, ok := in.(StgtH)
	require.True(t, ok, "type = %T, want StgtH", in)
	require.Equal(t, "stgt.h $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x387cb9ac), ctorWord(t, x))
}
