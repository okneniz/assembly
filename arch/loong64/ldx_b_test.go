package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdxBCtor(t *testing.T) {
	// llvm-mc-verified: ldx.b $t0, $t1, $t2.
	in := New().LdxB(lreg(t, 12), lreg(t, 13), lreg(t, 14))
	require.Equal(t, uint32(0x380039ac), ctorWord(t, in))

	_, ok := in.(LdxB)
	require.True(t, ok, "type = %T, want LdxB", in)
}

func TestLdxBDecodeEncode(t *testing.T) {
	in := decodeLdxB(0x380039ac, 0x90000000)

	x, ok := in.(LdxB)
	require.True(t, ok, "type = %T, want LdxB", in)
	require.Equal(t, "ldx.b $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x380039ac), ctorWord(t, x))
}
