package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdxWuCtor(t *testing.T) {
	// llvm-mc-verified: ldx.wu $t0, $t1, $t2.
	in := NewLdxWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))
	require.Equal(t, uint32(0x382839ac), ctorWord(t, in))

	_, ok := in.(LdxWu)
	require.True(t, ok, "type = %T, want LdxWu", in)
}

func TestLdxWuDecodeEncode(t *testing.T) {
	in := decodeLdxWu(0x382839ac, 0x90000000)

	x, ok := in.(LdxWu)
	require.True(t, ok, "type = %T, want LdxWu", in)
	require.Equal(t, "ldx.wu $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x382839ac), ctorWord(t, x))
}
