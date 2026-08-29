package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRotrWCtor(t *testing.T) {
	// llvm-mc-verified: rotr.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001b39ac),
		ctorWord(t, NewRotrW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewRotrW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(RotrW)
	require.True(t, ok, "type = %T, want RotrW", in)
}

func TestRotrWDecodeEncode(t *testing.T) {
	in := decodeOne(0x001b39ac, 0x90000000)

	x, ok := in.(RotrW)
	require.True(t, ok, "type = %T, want RotrW", in)
	require.Equal(t, "rotr.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x001b39ac), ctorWord(t, x))
}
