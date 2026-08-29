package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRotrDCtor(t *testing.T) {
	// llvm-mc-verified: rotr.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001bb9ac),
		ctorWord(t, NewRotrD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewRotrD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(RotrD)
	require.True(t, ok, "type = %T, want RotrD", in)
}

func TestRotrDDecodeEncode(t *testing.T) {
	in := decodeOne(0x001bb9ac, 0x90000000)

	x, ok := in.(RotrD)
	require.True(t, ok, "type = %T, want RotrD", in)
	require.Equal(t, "rotr.d $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x001bb9ac), ctorWord(t, x))
}
