package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdleBCtor(t *testing.T) {
	// llvm-mc-verified: ldle.b $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387a39ac),
		ctorWord(t, NewLdleB(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewLdleB(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(LdleB)
	require.True(t, ok, "type = %T, want LdleB", in)
}

func TestLdleBDecodeEncode(t *testing.T) {
	in := decodeLdleB(0x387a39ac, 0x90000000)

	x, ok := in.(LdleB)
	require.True(t, ok, "type = %T, want LdleB", in)
	require.Equal(t, "ldle.b $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x387a39ac), ctorWord(t, x))
}
