package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestXorCtor(t *testing.T) {
	// llvm-mc-verified: xor $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0015b9ac),
		ctorWord(t, New().Xor(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().Xor(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(Xor)
	require.True(t, ok, "type = %T, want Xor", in)
}

func TestXorDecodeEncode(t *testing.T) {
	in := decodeOne(0x0015b9ac, 0x90000000)

	x, ok := in.(Xor)
	require.True(t, ok, "type = %T, want Xor", in)
	require.Equal(t, "xor $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x0015b9ac), ctorWord(t, x))
}
