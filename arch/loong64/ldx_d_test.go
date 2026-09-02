package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdxDCtor(t *testing.T) {
	// llvm-mc-verified: ldx.d $t0, $t1, $t2.
	in := New().LdxD(lreg(t, 12), lreg(t, 13), lreg(t, 14))
	require.Equal(t, uint32(0x380c39ac), ctorWord(t, in))

	_, ok := in.(LdxD)
	require.True(t, ok, "type = %T, want LdxD", in)
}

func TestLdxDDecodeEncode(t *testing.T) {
	in := decodeLdxD(0x380c39ac, 0x90000000)

	x, ok := in.(LdxD)
	require.True(t, ok, "type = %T, want LdxD", in)
	require.Equal(t, "ldx.d $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x380c39ac), ctorWord(t, x))
}
