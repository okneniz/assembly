package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLddirCtor(t *testing.T) {
	v, err := New().UImm8(1)
	require.NoError(t, err)

	// llvm-mc-verified: lddir $t0, $t1, 1.
	require.Equal(
		t,
		uint32(0x064005ac),
		ctorWord(t, New().Lddir(lreg(t, 12), lreg(t, 13), v)),
	)

	in := New().Lddir(lreg(t, 12), lreg(t, 13), v)
	_, ok := in.(Lddir)
	require.True(t, ok, "type = %T, want Lddir", in)
}

func TestLddirDecodeEncode(t *testing.T) {
	// llvm-mc-verified: lddir $t0, $t1, 1.
	in := decodeLddir(0x064005ac, 0x90000000)

	x, ok := in.(Lddir)
	require.True(t, ok, "type = %T, want Lddir", in)
	require.Equal(t, "lddir $t0, $t1, 1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x064005ac), ctorWord(t, x))
}
