package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCloDCtor(t *testing.T) {
	// llvm-mc-verified: clo.d $t0, $t1.
	require.Equal(
		t,
		uint32(0x000021ac),
		ctorWord(t, New().CloD(lreg(t, 12), lreg(t, 13))),
	)

	in := New().CloD(lreg(t, 1), lreg(t, 2))
	_, ok := in.(CloD)
	require.True(t, ok, "type = %T, want CloD", in)
}

func TestCloDDecodeEncode(t *testing.T) {
	in := decodeCloD(0x000021ac, 0x90000000)

	clod, ok := in.(CloD)
	require.True(t, ok, "type = %T, want CloD", in)
	require.Equal(t, "clo.d $t0, $t1", clod.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), clod.Addr())
	require.Equal(t, uint32(0x000021ac), ctorWord(t, clod))
}
